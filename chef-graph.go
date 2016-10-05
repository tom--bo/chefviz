package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
    "github.com/awalterschulze/gographviz"
)

func main() {
    if len(os.Args) < 2 {
        panic("No role or recipes")
    }

    recipeRegistory := make(map[string]([]string))
    var recipes []string
    recipes = append(recipes, normalizeRecipeName(os.Args[1]))

    graphAst, _ := gographviz.Parse([]byte(`digraph G{}`))
    graph := gographviz.NewGraph()
    gographviz.Analyse(graphAst, graph)

    for len(recipes) > 0 {
        recipe := recipes[0]
        _, ok := recipeRegistory[recipe]
        if ok == true {
            // duplicate recipes
            recipes = recipes[1:]
            continue
        }
        recipepath := recipeToFilename(recipe)
        included, err := searchRecipesFromFile(recipepath)
        if err != nil {
            fmt.Println("`" + recipes[0] + "` is not found.")
            recipes = recipes[1:]
            continue
        }

        recipeRegistory[recipe] = included
        recipes = append(recipes, included...)
        addGraph(graph, recipes[0], included)
        recipes = recipes[1:]
    }
    fmt.Println(graph.String())
}

func normalizeRecipeName (recipe string) string {
    tmp := strings.Split(recipe, "::")
    if len(tmp) > 1 {
        return recipe
    }
    return recipe + "::default"
}

func recipeToFilename (recipe string) string {
    tmp := strings.Split(recipe, "::")
    return "../sample-chef-repo/cookbooks/" + tmp[0] + "/recipes/" + tmp[1] + ".rb"
}

func searchRecipesFromFile (filename string) ([]string, error){
    var fp *os.File
    var err error
    var ret []string

    fp, err = os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer fp.Close()

    re := regexp.MustCompile(".*include_recipe\\s+[\\'\"]([a-zA-Z0-9]+[:a-zA-Z0-9-_]*)[\\'\"].*")
    scanner := bufio.NewScanner(fp)
    for scanner.Scan() {
        line := scanner.Text()

        recipename:= re.FindStringSubmatch(line)
        if len(recipename) > 0 {
            ret = append(ret, normalizeRecipeName(recipename[1]))
        }
    }

    return ret, nil
}

func addGraph(graph *gographviz.Graph, parent string, children []string) {
    graph.AddNode("G", `"`+parent+`"`, nil)
    for _, child := range children {
        graph.AddNode("G", `"`+child+`"`, nil)
        graph.AddEdge(`"`+parent+`"`, `"`+child+`"`, true, nil)
    }
}
