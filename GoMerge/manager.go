package GoMerge

import (
	"github.com/bitly/go-simplejson"
	"fmt"
	"strings"
	"os"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"errors"
	"crypto/md5"
	"encoding/hex"
)

type Manager struct {
	Net         *Net
	UpTagText   string
	IdBase      int
	ProjectList map[string]*Project
	WatcherList map[string]*Watcher
}

func (p *Manager)Init() {
	p.UpTagText = "/**\n* [nowTime]\n* caoping@163.com\n**/\n\n"
	p.ProjectList = map[string]*Project{}
	p.WatcherList = map[string]*Watcher{}
	p.Net = new(Net)
	p.Net.Socket = "6301"
	p.Net.manager = p
	p.Load()
	p.Net.Start()
	for _, v := range p.WatcherList {
		v.manager = p
		v.Run()
	}
}

func (p *Manager)addWatchFile(m *simplejson.Json) {
	fmt.Printf("addWatchFile fullFile:\n%v\n", m.Get("data").MustString())
	project := new(Project)
	project.UpTagText = p.UpTagText
	if e := project.Init(m.Get("data").MustString()); e != nil {
		fmt.Printf("addWatchFile:%v\n", e)
		p.Net.Send(Msg{Cmd:"addWatchFile", Status:1, Msg:"项目添加失败!"})
		return
	}
	if p.check(project.Hex) {
		fmt.Printf("此项目添加过,不需要重复添加!\n")
		p.Net.Send(Msg{Cmd:"addWatchFile", Status:1, Msg:"此项目添加过,不需要重复添加!"})
		return
	}
	project.Id = p.IdBase
	p.IdBase++
	p.ProjectList[project.Hex] = project
	p.watchProject([]string{strconv.Itoa(project.Id)})
	p.Net.Send(Msg{Cmd:"addWatchFile", Msg:"项目添加成功!"})
	p.getProjectList()
	p.Save()
}

func (p *Manager)check(hex string) bool {
	for _, v := range p.ProjectList {
		if v.Hex == hex {
			return true
		}
	}
	return false
}

func (p *Manager)rount(s string) bool {
	m, err := simplejson.NewJson([]byte(s))
	if err != nil {
		return true
	}
	cmd := m.Get("cmd").MustString()
	switch cmd{
	case "addWatchFile":
		p.addWatchFile(m)
	case "getProjectList":
		p.getProjectList()
	case "getWatcherList":
		p.getWatcherList()
	case "update":
		p.Update([]string{m.Get("hex").MustString()})
	case "rebuild":
		p.Rebuild([]string{m.Get("hex").MustString()})
	case "delete":
		p.Delete([]string{m.Get("hex").MustString()})
	default:

	}
	return false
}

func (p *Manager)CmdRout(t string) {
	if len(t) == 0 {
		fmt.Printf("\n")
		return
	}
	a := strings.Fields(t)
	switch a[0]{
	case "exit", "x":
		fmt.Printf("退出!\n")
		os.Exit(0)
	case "list":
		p.list(a)
	case "save":
		p.Save()
	case "load":
		p.Load()
	case "watchProject":
		p.watchProject(a[1:])
	case "update":
		p.Update(a[1:])
	default:
		fmt.Println(t)
	}
}

func (p *Manager)list(a []string) {
	for _, v := range p.ProjectList {
		isWatch := "x"
		if v.NoWatch {
			isWatch = " "
		}
		isUpdate := " "
		if v.IsUpdate {
			isUpdate = "x"
		}
		fmt.Printf("[ %v ][%s][%s] %v\n", v.Id, isWatch, isUpdate, v.FileName)
	}

}

func (t *Manager)watchProject(a []string) {
	project, err := t.getProject(a)
	if err != nil {
		fmt.Printf("监控项目失败!\n%s\n", err)
		return
	}
	for _, v := range project.MtoList {
		for _, src := range v.Src {
			t.addToWatch(project.Hex, src)
		}
	}
	t.getWatcherList()
	t.Save()
}

func (p *Manager)addToWatch(hex, src string) {
	lastPoint := strings.LastIndex(src, "\\")
	lastPoint2 := strings.LastIndex(src, "/")
	if lastPoint < lastPoint2 {
		lastPoint = lastPoint2
	}
	root := src[:lastPoint + 1]
	name := src[lastPoint + 1:]
	w, ok := p.WatcherList[root]
	if ok == true {
		w.add(hex, name)
		return
	}
	watch := new(Watcher)
	watch.Root = root
	watch.init()
	watch.add(hex, name)
	watch.manager = p
	watch.Run()
	p.WatcherList[root] = watch
}

func (t *Manager)getProject(a []string) (project *Project, e error) {
	if len(a) == 0 {
		e = errors.New("请输入执行参数\n")
		return project, e
	}
	id, e := strconv.Atoi(a[0])
	if e == nil {
		project, e = t.getProjectById(id)
	} else {
		project, e = t.getProjectByHex(a[0])
	}
	return project, e

}

func (p *Manager)getProjectById(id int) (*Project, error) {
	for _, v := range p.ProjectList {
		if v.Id == id {
			return v, nil
		}
	}
	return nil, errors.New("未找到对应的项目!")
}
func (p *Manager)getProjectByHex(hex string) (*Project, error) {
	for _, v := range p.ProjectList {
		if v.Hex == hex {
			return v, nil
		}
	}
	return nil, errors.New("未找到对应的项目!")
}
//获取项目列表
func (p *Manager)getProjectList() {
	data, e := json.Marshal(p.ProjectList)
	if e != nil {
		p.Net.Send(Msg{Cmd:"getProjectList", Status:1, Msg:"获取项目列表失败!"})
		return
	}
	p.Net.Send(Msg{
		Cmd:"getProjectList",
		Status:0,
		Data:string(data),
	})
}
//获取监控列表
func (p *Manager)getWatcherList() {
	data, e := json.Marshal(p.WatcherList)
	if e != nil {
		p.Net.Send(Msg{Cmd:"getWatcherList", Status:1, Msg:"获取监控列表失败!"})
		return
	}
	p.Net.Send(Msg{
		Cmd:"getWatcherList",
		Status:0,
		Data:string(data),
	})
}

func (t *Manager)Update(a []string) {
	if len(a) == 0 {
		fmt.Printf("请输入执行参数\n")
		return
	}
	id, e := strconv.Atoi(a[0])
	var project *Project
	var err error
	if e == nil {
		project, err = t.getProjectById(id)
	} else {
		project, err = t.getProjectByHex(a[0])
	}
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	project.Merge()
	t.Net.Send(Msg{
		Cmd:"onUpdate",
		Msg:"项目更新 " + project.FileName,
		Data:map[string]string{
			"ProjectHex":project.Hex,
		},
	})
	fmt.Printf("更新!\n")
}

func (t *Manager)Clear(projectHex string) {
	for _, v := range t.WatcherList {
		v.clear(projectHex)
	}

	for k, v := range t.WatcherList {
		if len(v.FileList) == 0 {
			delete(t.WatcherList,k)
		}
	}
}

/*
项目删除
*/
func (t *Manager)Delete(a []string) {
	project, err := t.getProject(a)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	t.Clear(project.Hex)
	delete(t.ProjectList,project.Hex)
	msg := "项目删除 " + project.FileName + "\n"
	t.Net.Send(Msg{
		Cmd:"onDelete",
		Msg:msg,
		Data:map[string]string{},
	})
	t.Save()
	t.RefreshWeb()
	fmt.Printf(msg)
}
/*
项目重建
*/
func (t *Manager)Rebuild(a []string) {
	project, err := t.getProject(a)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	t.Clear(project.Hex)
	project.Init(project.FileName)
	t.watchProject([]string{project.Hex})
	msg := "项目重建 " + project.FileName + "\n"
	t.Net.Send(Msg{
		Cmd:"onRebuild",
		Msg:msg,
		Data:map[string]string{
			"ProjectHex":project.Hex,
		},
	})

	t.RefreshWeb()
	fmt.Printf(msg)
}

func (t *Manager)RefreshWeb() {
	t.getProjectList()
	t.getWatcherList()
}

func (p *Manager)Save() {
	j, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("www/conf.json", []byte(j), 0666)
	if err != nil {
		panic(err)
	}
	fmt.Printf("保存参数完毕!\n")
}

func (p *Manager)Load() {
	t, e := ioutil.ReadFile("www/conf.json")
	if e != nil {
		fmt.Printf("未找到配置文件!\n")
		return
	}
	e = json.Unmarshal(t, &p)
	if e != nil {
		fmt.Printf("参数编码失败!\n")
	} else {
		fmt.Printf("读取参数完毕!\n")
	}
}

/*
获取md5字符
*/
func getHex(str string) string {
	md := md5.New()
	md.Write([]byte(str))
	s := md.Sum(nil)
	return hex.EncodeToString(s)
}