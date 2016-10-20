package core

import (
	"path/filepath"
	"github.com/pandazhuzi/buns/utils"
	"github.com/pandazhuzi/buns/errors"
	"fmt"
	"os/exec"
	"os"
	"runtime"
)

type CliProject struct {
	Name string
	root string
	source string
	build string
	mainGoFile string
	cmds string
	program string
}


func CreateCliProject(name string, resource string, targer string) (*CliProject,error) {

	targer = filepath.Join(targer, name)
	archive := filepath.Join(resource, "archive")


	err := utils.FolderCopy(archive, targer)

	if(err != nil){
		return nil,errors.MakeError(err)
	}

	cli, err := OpenCliProject(targer)

	if(err != nil){
		return nil,err
	}


	err = utils.FolderRenameByTemplate(targer, cli)

	if(err != nil){
		return nil,err
	}

	err = utils.FolderRenderByTemplate(targer, ".tpl." ,cli)

	if(err != nil){
		return nil,errors.MakeError(err)
	}

	return cli,nil
}

func OpenCliProject(root string) (*CliProject, error){

	if(!utils.FileExists(root)){
		return nil,errors.MakeError("open cli project error , folder %v not exist", root)
	}

	cli := new(CliProject)

	cli.Name = filepath.Base(root)
	cli.root = root
	cli.source = filepath.Join(cli.root,"source")
	cli.build = filepath.Join(cli.root, "build")
	cli.mainGoFile = filepath.Join(cli.source,cli.Name + ".go")
	cli.cmds = filepath.Join(cli.source,"cmd")
	cli.program = filepath.Join(cli.build,cli.Name)

	if(runtime.GOOS == "windows"){
		cli.program += ".exe"
	}

	return cli,nil

}



func (p *CliProject) Add(name string, resource string, targer string) error{

	if(len(name) == 0){
		return errors.MakeError("name is empty.")
	}


	type addOptions struct {
		Project *CliProject
		Name string
		CmdFolder string
		CamelName string
		UnixName string
	}

	opts := new(addOptions)
	opts.CmdFolder = p.cmds
	opts.Name = name
	opts.CamelName,opts.UnixName = formatName(opts.Name)

	resource = filepath.Join(resource,"templates","add.tpl.go")
	targer = filepath.Join(targer,"source","cmd",fmt.Sprintf("%v.tpl.go", opts.Name))

	err := utils.FileCopy(resource,targer)

	if(err != nil){
		return err
	}

	err = utils.RenderTemplateFile(targer,".tpl.",opts)

	if(err != nil){
		return err
	}

	return nil

}

func (p *CliProject) runCommandLine(name string, args ...string) error{

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = p.root

	err := cmd.Run()

	if(err != nil){
		return errors.MakeError(err)
	}

	return nil

}

func (p *CliProject) Build() error{

	err := p.runCommandLine(
		"go","build",
		"-o",p.program,
		"-v",
		p.mainGoFile)

	if(err != nil){
		return errors.MakeError(err)
	}

	return nil

}

func (p *CliProject) Package(targer string) error {

	err := p.Build()

	if(err != nil){
		return err
	}

	info, err := os.Stat(targer)

	if(err != nil){
		return err
	}


	if(info.IsDir()){
		targer = filepath.Join(targer, p.Name + ".tar.gz")
	}

	err = utils.Tar(p.build,targer,p.Name)

	if(err != nil){
		return errors.MakeError(err)
	}

	return nil
}