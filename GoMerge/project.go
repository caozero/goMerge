package GoMerge

import (
	"io/ioutil"
	"strings"
	"regexp"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type Mto struct {
	Hex string `json:"hex"`
	Src []string `json:"src"`
	MergeTo string `json:"mergeTo"`
}

func (t *Mto)merge(UpTagText string){
	outText:=UpTagText
	for _,v:=range t.Src{
		fmt.Printf("读取:    %s\n", v)
		read, err := ioutil.ReadFile(v)
		if err != nil {
			continue
		}
		outText+=string(read)
	}
	err:=ioutil.WriteFile(t.MergeTo,[]byte(outText),0666)
	if err!=nil{
		fmt.Printf("[  %s ] 合并文件写入失败!\n", t.MergeTo)
		return
	}
	fmt.Printf("写入:    %s\n", t.MergeTo)
}

type Project struct {
	Id int
	Name string
	Hex string
	FileName string
	RootPath string
	IsUpdate bool
	NoWatch bool
	NotUpdate bool
	UpdateTime time.Time
	UpTagText string
	MtoList      map[string]*Mto `json:"mto"`
}

func (p *Project)Init(fileName string)error{
	outText, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	p.MtoList=map[string]*Mto{}
	s2:=p.clearNotes(string(outText))
	md := md5.New()
	md.Write([]byte(fileName))
	s := md.Sum(nil)
	smh := hex.EncodeToString(s)
	p.Hex = smh
	p.FileName = fileName
	p.RootPath = fileName[0:strings.LastIndex(fileName, "\\") + 1]
	p.MustCompile([]byte(s2))
	return err
}

/*
清理掉注释
*/
func (p *Project)clearNotes(s string)string{
	start:=strings.Index(s,"<!--")
	if start==-1{
		return s
	}
	s2:=s[start+4:]
	end:=strings.Index(s2,"-->")
	if end==-1{
		return s
	}
	s=s[:start]+s2[end+3:]
	return p.clearNotes(s)
}

func (p *Project)MustCompile(s []byte){
	reg := regexp.MustCompile(`<script[^>]*?>.*?</script>`)
	a := reg.FindAll(s, -1)
	for _, r := range a {
		if strings.Index(string(r), "mergeTo") != -1 {
			s := string(r)
			p.GetInLine(s)
		}

	}
	reg = regexp.MustCompile(`<link[^>]*?>`)
	css := reg.FindAll(s, -1)
	for _, r := range css {
		if strings.Index(string(r), "mergeTo") != -1 {
			s := string(r)
			p.GetInLine(s)
		}
	}
	return
}

func (m *Project)GetInLine(s string) {
	a := strings.Fields(s)
	src:=""
	mergeTo:=""
	for _, p := range a {
		if strings.Index(p, "=") != -1 {
			end := strings.Index(p, ">");
			if end != -1 {
				p = p[0:end]
			}
			params := strings.Split(p, "=")
			if len(params) > 1 {
				if params[0] == "src" {
					src=getRootPath(m.RootPath,strings.Trim(params[1], "\""));
				}
				if params[0] == "href" {
					src=getRootPath(m.RootPath,strings.Trim(params[1], "\""));
				}
				if params[0] == "mergeTo" {
					mergeTo=getRootPath(m.RootPath,strings.Trim(params[1], "\"/"));
				}

			}
		}
	}

	md := md5.New()
	md.Write([]byte(mergeTo))
	hex:=hex.EncodeToString(md.Sum(nil))
	fmt.Printf("hex:%v\n",hex)
	if mto,ok:=m.MtoList[hex];ok==true{
		mto.Src=append(mto.Src,src)
		fmt.Printf("添加相同目标:%v\n",src)
		return
	}
	fmt.Printf("新建立目标:%v\n",mergeTo)
	m.MtoList[hex]=&Mto{
		Hex:hex,
		Src:[]string{src},
		MergeTo:mergeTo,
	}
	return
}

func (t *Project)Merge() {
	if t.NotUpdate{
		fmt.Printf("设定为不参与更新,  %s\n",t.FileName)
		return
	}
	fmt.Printf("合并更新操作:\n %s\n", t.FileName)
	s:=formatText(t.UpTagText)
	for _, v := range t.MtoList {
		v.merge(s)
	}
	t.UpdateTime=time.Now()
	t.IsUpdate=false
	fmt.Printf("完毕!\n")
}

//格式化文字
func formatText(s string)string{
	s=strings.Replace(s,"[nowTime]",string(time.Now().Format("2006-01-02 15:04:05")),-1)
	return s
}

func getRootPath(root, inPath string) string {
	inPath = strings.Replace(inPath, "/", "\\", -1)
	dot := strings.Count(inPath, "..\\")
	inPath = inPath[dot * 3:]
	root = root[:len(root) - 1]
	for dot > 0 {
		root = root[:strings.LastIndex(root, "\\")]
		dot--
	}
	fmt.Printf("root:%v\n",root)
	fmt.Printf("inPath:%v\n",inPath)
	rootFilePath:=""
	if string(root[len(root)-1])=="\\" || string(inPath[0])=="\\"{
		rootFilePath = root + inPath
	}else{
		rootFilePath = root + "\\" + inPath
	}

	return rootFilePath
}