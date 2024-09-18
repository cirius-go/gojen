package gojen_test

// func TestStore(t *testing.T) {
// 	t.Run("It should add a decl to the store", func(t *testing.T) {
// 		cliCfg := gojen.ConsoleC()
// 		c := gojen.NewConsoleWithConfig(cliCfg)
// 		fileManagerCfg := gojen.FileManagerC()
// 		f := gojen.NewFileManagerWithConfig(fileManagerCfg)
// 		cfg := gojen.StoreC()
// 		s := gojen.NewStoreWithConfig(cfg, c, f)
// 		d := &gojen.D{
// 			Path:         "custom/path",
// 			Require:      []string{},
// 			Name:         "decl1",
// 			Args:         map[string]any{},
// 			Elements:     []*gojen.E{},
// 			Dependencies: []string{},
// 			Description:  "Test decl",
// 		}
// 		s.SetDecl(d)
//
// 		assert.Equal(t, d, s.GetDecl("decl1"))
// 	})
// }
//
// func TestLoadDir(t *testing.T) {
// 	t.Run("It should load all decls from yaml files", func(t *testing.T) {
// 		dirPath := testlib.CreateDir(t, ".gojen/decls")
// 		testlib.NewFileWithContent(t, filepath.Join(dirPath, "decl1.yaml"),
// 			`
// path: "example/go/internal/service/{{ sLower .Domain }}.go"
// require: ["Domain"]
// name: "service"
// args:
//   domain: "user"
// elements:
//   - name: "initFile"
//     template: "package service"
//     strategy: "init"
//     confirm: true
// `)
//
// 		cliCfg := gojen.ConsoleC()
// 		c := gojen.NewConsoleWithConfig(cliCfg)
//
// 		fileManagerCfg := gojen.FileManagerC()
// 		f := gojen.NewFileManagerWithConfig(fileManagerCfg)
//
// 		cfg := gojen.StoreC()
// 		s := gojen.NewStoreWithConfig(cfg, c, f)
//
// 		err := s.LoadDir(dirPath)
// 		assert.Nil(t, err)
//
// 		d := s.GetDecl("service")
// 		assert.Equal(t, "example/go/internal/service/{{ sLower .Domain }}.go", d.Path)
// 		assert.Equal(t, []string{"Domain"}, d.Require)
// 		assert.Equal(t, len(d.Elements), 1)
// 		assert.Equal(t, "service", d.Name)
// 	})
//
// 	t.Run("It should ask user to replace decl if it already exists in the store", func(t *testing.T) {
// 		dirPath := testlib.CreateDir(t, ".gojen/decls")
// 		testlib.NewFileWithContent(t, filepath.Join(dirPath, "decl1.yaml"),
// 			`
// path: "example/go/internal/service/{{ sLower .Domain }}.go"
// require: ["Domain"]
// name: "service"
// args:
//   domain: "user"
// elements:
//   - name: "initFile"
//     template: "package service"
//     strategy: "init"
//     confirm: true
// `)
//
// 		testlib.NewFileWithContent(t, filepath.Join(dirPath, "decl2.yaml"),
// 			`
// path: "example/go/internal/service2/{{ sLower .Domain }}.go"
// require: ["Domain"]
// name: "service"
// args:
//   domain: "aws"
// elements:
//   - name: "initFile2"
//     template: "package service2"
//     strategy: "init"
//     confirm: true
// `)
//
// 		c := mocks.NewConsoleManager(t)
// 		c.Mock.On("Infof", mock.AnythingOfType("bool"), mock.AnythingOfType("string"), mock.Anything).Return(true)
// 		c.On("PerformYesNo", "Declaration with name '%s' already exists. Do you want to override it?\n", "service").Return(true)
//
// 		fileManagerCfg := gojen.FileManagerC()
// 		f := gojen.NewFileManagerWithConfig(fileManagerCfg)
//
// 		cfg := gojen.StoreC()
// 		s := gojen.NewStoreWithConfig(cfg, c, f)
//
// 		err := s.LoadDir(dirPath)
// 		assert.Nil(t, err)
//
// 		d := s.GetDecl("service")
// 		assert.Equal(t, d.Args["domain"], "aws")
// 		assert.Equal(t, "example/go/internal/service2/{{ sLower .Domain }}.go", d.Path)
// 		assert.Equal(t, []string{"Domain"}, d.Require)
// 		assert.Equal(t, len(d.Elements), 1)
// 		assert.Equal(t, "service", d.Name)
// 	})
// }
