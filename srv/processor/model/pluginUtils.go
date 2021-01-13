package model

import (
	"fmt"
	"github.com/sidazhang123/f10-go/srv/processor/plugins/debug"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
)

func (s *Service) GetPluginPath() (string, error) {
	pluginSrcPath := Params.PluginSrcPath
	files, err := ioutil.ReadDir(pluginSrcPath)
	log.Info("pluginSrcPath=" + pluginSrcPath)
	if err != nil {
		log.Error(err.Error())
		return "nil", err
	}
	res := ""
	for _, f := range files {
		dirPath, err := filepath.Abs(pluginSrcPath)
		if err != nil {
			log.Error(err.Error())
			return "nil", err
		}
		res += filepath.Join(dirPath, f.Name()) + ";"
	}
	return strings.Trim(res, ";"), nil
}
func (s *Service) GetPluginSrc(path string) (string, error) {
	b, err := ioutil.ReadFile(path) // b has type []byte
	if err != nil {
		log.Error(err.Error())
	}
	return string(b), nil
}

func (s *Service) BuildSo(newSrcCode, path string) (string, error) {
	//update source code
	err := ioutil.WriteFile(path, []byte(newSrcCode), 0644)
	if err != nil {
		return "nil", fmt.Errorf("failed to update .go\n%s", err.Error())
	}
	//plugin must be in main package
	err = togglePackageName(path)
	if err != nil {
		return "nil", fmt.Errorf("failed to togglePackageName\n%s", err.Error())
	}
	// make go build output path and filename
	soPath, err := filepath.Abs(Params.PluginSoPath)
	if err != nil {
		return "nil", fmt.Errorf("failed to get Abs(Params.PluginSoPath)\n%s", err.Error())
	}
	soPath = filepath.Join(soPath, strings.TrimSuffix(filepath.Base(path), "go")+"so")

	// exec the build cmd
	var output []byte
	output, err = exec.Command("go", "build", "-buildmode=plugin", "-o", soPath, path).CombinedOutput()
	if err != nil {
		_ = togglePackageName(path)
		return "nil", fmt.Errorf("failed to build\n%s\n%s", err.Error(), string(output))
	}
	// change package name back to src
	err = togglePackageName(path)
	if err != nil {
		return "nil", fmt.Errorf("failed to togglePackageName\n%s", err.Error())
	}
	// load the new so
	pdll, err := plugin.Open(soPath)
	if err != nil {
		return "nil", fmt.Errorf("reload So::failed to open %s\n%s", soPath, err.Error())
	}
	funcName := strings.TrimSuffix(filepath.Base(soPath), ".so")
	//func name in plugin must begin with a capital letter
	f, err := pdll.Lookup(strings.Title(funcName))
	if err != nil {
		return "nil", fmt.Errorf("reload So::failed to find func name %s\n%s", strings.Title(funcName), err.Error())
	}

	s.RegexFunc[funcName] = f
	return soPath, nil
}

func togglePackageName(path string) error {
	pack := map[string]string{"src": "main", "main": "src"}

	input, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	success := false
	for i, line := range lines {
		if success {
			break
		}
		for k, v := range pack {
			if strings.HasPrefix(strings.TrimSpace(line), "package "+k) {
				lines[i] = "package " + v
				success = true
				break
			}
		}
	}
	if !success {
		return fmt.Errorf("package name line not found")
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RegexTest(testStr, pluginPath string) (map[string]string, error) {
	//log.Info(fmt.Sprintf("RegexTest pluginPath: %+v",pluginPath))
	funcName := strings.TrimSuffix(filepath.Base(pluginPath), ".go")
	//log.Info(fmt.Sprintf("RegexTest funcName: %+v",funcName))
	resMap, err := s.CallSoPlugin(funcName, testStr)
	//log.Info(fmt.Sprintf("RegexTest resMap: %+v",resMap))
	//log.Info(fmt.Sprintf("RegexTest err: %+v",err))
	if len(err) > 0 {
		eStr := ""
		for _, e := range err {
			eStr += e.Error()
		}
		return nil, fmt.Errorf(eStr)
	}
	r := map[string]string{}
	for k, v := range resMap {
		r[k] = fmt.Sprint(v)
	}
	return r, nil
}

func (s *Service) RegisterPlugin(level string) error {
	if level == "prod" {
		return registerSo(s)
	} else {
		return registerGo(s)
	}
}
func registerGo(s *Service) error {
	files, err := ioutil.ReadDir(Params.PluginGoPath)
	if err != nil {
		return err
	}
	var d debug.Debug
	exclFuncName := strings.Split(Params.PluginExcl, ";")
	for _, file := range files {
		funcName := strings.TrimSuffix(filepath.Base(file.Name()), ".go")
		// skip excl func noted in env.yml
		if true == func(string) bool {
			if funcName == "struct" {
				return true
			}
			for _, i := range exclFuncName {
				if i == funcName {
					return true
				}
			}
			return false
		}(funcName) {
			continue
		}
		m := reflect.ValueOf(&d).MethodByName(strings.Title(funcName))
		if !m.IsValid() {
			return fmt.Errorf("reflect can't find funcName: %s", funcName)
		}
		s.RegexFunc[funcName] = m

	}
	return nil
}

func (s *Service) CallGoPlugin(funcName, testStr string) (map[string]interface{}, []error) {
	if v, ok := s.RegexFunc[funcName]; ok {
		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(testStr)
		res := v.(reflect.Value).Call(params)
		return res[0].Interface().(map[string]interface{}), res[1].Interface().([]error)
	} else {
		return nil, []error{fmt.Errorf("func %s not registered", funcName)}
	}
}

func registerSo(s *Service) error {
	files, err := ioutil.ReadDir(Params.PluginSoPath)
	if err != nil {
		return err
	}
	exclFuncName := strings.Split(Params.PluginExcl, ";")
	for _, file := range files {
		funcName := strings.TrimSuffix(filepath.Base(file.Name()), ".so")
		// skip excl func noted in env.yml
		if true == func(string) bool {
			for _, i := range exclFuncName {
				if i == funcName {
					return true
				}
			}
			return false
		}(funcName) {
			continue
		}

		dirPath, err := filepath.Abs(Params.PluginSoPath)
		if err != nil {
			return err
		}
		soPath := filepath.Join(dirPath, file.Name())
		pdll, err := plugin.Open(soPath)
		if err != nil {
			return fmt.Errorf("failed to open %s\n%s", soPath, err.Error())
		}
		//func name in plugin must begin with a capital letter
		f, err := pdll.Lookup(strings.Title(funcName))
		if err != nil {
			return fmt.Errorf("failed to find func name %s\n%s", strings.Title(funcName), err.Error())
		}

		s.RegexFunc[funcName] = f

	}
	//log.Info(fmt.Sprintf("registerSo %+v",s.RegexFunc))
	return nil
}

// .so name(without ext) == pluginName == funcName == flagName
// return {code,flagname,updatetime,fields...} from plugin
func (s *Service) CallSoPlugin(funcName, testStr string) (map[string]interface{}, []error) {
	//log.Info(fmt.Sprintf("CallSoPlugin funcName %+v",funcName))
	f := s.RegexFunc[funcName].(func(string) (map[string]interface{}, []error))
	mapToDb, errList := f(testStr)
	if errList != nil && len(errList) > 0 {
		return nil, errList
	}
	return mapToDb, nil
}
