package GoMerge

import (
	"fmt"
	"log"
	"github.com/howeyc/fsnotify"
	"strings"
	"time"
)

type WatchFile struct {
	Hex string
	Name string
	NoWatch bool
	IsUpdate bool
	ModifyTime time.Time
	ProjectHex []string
}

func (t *WatchFile)add(projectHex string){
	for _,v:=range t.ProjectHex{
		if v==projectHex{
			fmt.Printf("该项目的此文件已经添加至监控.\nprojectHex:%s\nfileName:%s\n",projectHex,t.Name)
			return
		}
	}
	t.ProjectHex=append(t.ProjectHex,projectHex)
	fmt.Printf("添加至监控:\nprojectHex:%s\nfileName:%s\n",projectHex,t.Name)
}
func (t *WatchFile)setProjectUpdate(manager *Manager){
	for _,v:=range t.ProjectHex{
		p,ok:=manager.ProjectList[v]
		if ok==true{
			p.IsUpdate=true
		}
	}

}

func (t *WatchFile)clear(projectHex string){
	for k,v:=range t.ProjectHex{
		if v==projectHex{
			if len(t.ProjectHex)==1{
				t.ProjectHex=[]string{}
			}else{
				t.ProjectHex=append(t.ProjectHex[:k],t.ProjectHex[k+1:]...)
			}
			return
		}
	}

}

type Watcher struct {
	Hex string
	Root string
	FileList map[string]*WatchFile
	NoWatch bool
	manager *Manager
}

func (p *Watcher)init(){
	p.Hex = getHex(p.Root)
	p.FileList=map[string]*WatchFile{}

}
func (p *Watcher)add(hex,name string){
	wf,ok:=p.FileList[name]
	if ok==true{
		wf.add(hex)
		return
	}
	f:=&WatchFile{}
	f.add(hex)
	f.Name=name
	f.Hex=getHex(p.Root+name)
	p.FileList[name]=f
	fmt.Printf("添加监视文件:%s\n",name)
}

func (p *Watcher)List(){
	if len(p.FileList)==0{
		fmt.Printf("监控列表为空!\n")
		return
	}
	fmt.Printf("监视目录:\n%s\n",p.Root)
	for k,v:=range p.FileList{
		isWatch:="x"
		if v.NoWatch{
			isWatch=" "
		}
		isUpdate:=" "
		if v.IsUpdate{
			isUpdate="x"
		}
		fmt.Printf("[ %d ][%s][%s] %s ---> %s\n",k,isWatch,isUpdate,v.Name,v.ModifyTime.Format("2006-01-02 15:04:05"))
	}
}

func (p *Watcher)Run(){
	if len(p.Root)==0 {
		fmt.Printf("监视目录Root为空!,\n")
		return
	}
	fmt.Printf("watch root: %v\n",p.Root)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case ev := <-watcher.Event:

				if ev.IsCreate() {
					//log.Println("文件事件 新建:", ev)
				} else if ev.IsDelete() {
					//log.Println("文件事件 删除:", ev)
				} else if ev.IsModify() {
					//log.Println("文件事件 修改:", ev)
					p.OnModify(ev.Name)
				} else if ev.IsRename() {
					//log.Println("文件事件 重命名:", ev)
				} else if ev.IsAttrib() {
					//log.Println("文件事件 修改元数据:", ev)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Watch(p.Root)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("启动监视目录 %s\n",p.Root)
}

type FileModifyData struct {
	Root string
	Name string
}

func (p *Watcher)OnModify(name string){
	name=name[strings.LastIndex(name,"\\")+1:]
	for _,v:=range p.FileList{
		if v.Name==name{
			v.IsUpdate=true
			v.ModifyTime=time.Now()
			v.setProjectUpdate(p.manager)
			fmt.Printf("文件修改: %s%s \n",p.Root,name)
			p.manager.Net.Send(Msg{
				Cmd:"onFileModify",
				Msg:p.Root+name+" 修改",
				Data:map[string]string{
					"WatcherHex":p.Hex,
					"Root":p.Root,
					"Hex":v.Hex,
					"Name":name,
				},
			})
		}
	}
	p.manager.Save()
}

func (t *Watcher)clear(projectHex string){
	for _,v:=range t.FileList{
		v.clear(projectHex)
	}
	for k,v:=range t.FileList{
		if len(v.ProjectHex)==0{
			delete(t.FileList,k)
		}
	}
}