package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
    "github.com/awalterschulze/gographviz"
)

type recipes struct {
    role string
    file string
    includeFile []string
}


func main() {
    //fmt.Printf("[%q]", strings.Trim(" !!! Achtung! Achtung! !!! ", "! "))
    if len(os.Args) < 2 {
        panic("no role")
    }
    graphAst, _ := gographviz.Parse([]byte(`digraph G{}`))
    graph := gographviz.NewGraph()
    gographviz.Analyse(graphAst, graph)

    var roleRegistory map[string]([]string)
    fmt.Println(roleRegistory)

    var roles []string
    roles = append(roles, os.Args[1])

    for len(roles) > 0 {
        rolepath := roleToFilename(roles[0])
        included, err := searchRolesFromFile(rolepath)
        if err != nil {
            fmt.Println("`" + roles[0] + "` is not found.")
            roles = roles[1:]
            continue
        }
        roles = append(roles, included...)
        addGraph(graph, roles[0], included)
        roles = roles[1:]
    }
    //files := search_roles_from_file(os.Args[1])
    fmt.Println(graph.String())
}

func roleToFilename (role string) string {
    tmp := strings.Split(role, "::")
    if len(tmp) > 1 {
        return "../sample-chef-repo/cookbooks/" + tmp[0] + "/recipes/" + tmp[1] + ".rb"
    }
    return "../sample-chef-repo/cookbooks/" + tmp[0] + "/recipes/default.rb"
}

func searchRolesFromFile (filename string) ([]string, error){
    var fp *os.File
    var err error
    var ret []string

    fp, err = os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer fp.Close()

    re := regexp.MustCompile("['\"].*?['\"]")
    scanner := bufio.NewScanner(fp)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Index(line, "include") != -1 {
            // fmt.Println(line)
            rolename:= strings.Trim(re.FindStringSubmatch(line)[0], "\"'")

            ret = append(ret, rolename)
        }
    }

    return ret, nil
}

func addGraph(graph *gographviz.Graph, parent string, children []string) {
    graph.AddNode("G", "\""+parent+"\"", nil)
    for _, child := range children {
        graph.AddNode("G", "\""+child+"\"", nil)
        graph.AddEdge("\""+parent+"\"", "\""+child+"\"", true, nil)
    }
}

