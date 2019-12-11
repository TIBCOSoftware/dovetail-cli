package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func CreateDirIfNotExist(subdir ...string) string {
	dir := path.Join(subdir...)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("createDirIfNotExist err %v", err)
			panic(err)
		}
	}

	return dir
}

func CopyFile(src string, dest string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return CopyContent(content, dest)
}
func CopyContent(content []byte, dest string) error {
	ft, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer ft.Close()
	ft.Write(content)
	return nil
}

func CreateTargetDirs(targetPath string) string {
	os.RemoveAll(targetPath)
	target := CreateDirIfNotExist(targetPath)
	return target
}

func CreateNewPom(oldpom []byte, targetdir, externalDepFile, newpomFileName string) error {
	dep := ""
	fmt.Printf("dep=%s, target=%s, pom=%s\n", externalDepFile, targetdir, newpomFileName)
	if externalDepFile != "" {
		deppom, err := ioutil.ReadFile(externalDepFile)
		if err != nil {
			return err
		}

		dep = string(deppom)
	}
	newpom := strings.Replace(string(oldpom), "%%external%%", dep, 1)

	err := CopyContent([]byte(newpom), path.Join(targetdir, newpomFileName))
	if err != nil {
		return err
	}
	return nil
}
func MvnPackage(groupId, artifactId, version, pomf, targetdir string) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	args := []string{"package", "-f", path.Join(targetdir, pomf), "-DbaseDir=" + targetdir, "-Dversion=" + version, "-DgroupId=" + groupId, "-DartifactId=" + artifactId}
	cmd := exec.Command("mvn", args...)
	logger.Printf("mvn command %v\n", cmd.Args)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("MvnPackage compile err %v", string(out))
	}
	fmt.Println(string(out))
	return nil
}

func MvnInstall(groupId, artifactId, version, file string) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	args := []string{"org.apache.maven.plugins:maven-install-plugin:2.5.2:install-file", "-Dpackaging=jar", "-DgeneratePom=true", "-Dversion=" + version, "-DgroupId=" + groupId, "-DartifactId=" + artifactId, "-Dfile=" + file}
	cmd := exec.Command("mvn", args...)
	logger.Printf("mvn command %v\n", cmd.Args)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("MvnInstall err %v", string(out))
	}
	fmt.Println(string(out))
	return nil
}
