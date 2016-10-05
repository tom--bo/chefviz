package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "regexp"
    "strings"
    "github.com/awalterschulze/gographviz"
)

type Chefviz struct {
	outStream, errStream io.Writer
    graph *gographviz.Graph
}

func (cv *Chefviz) newChefviz() {
    graphAst, _ := gographviz.Parse([]byte(`digraph G{}`))
    cv.graph = gographviz.NewGraph()
    gographviz.Analyse(graphAst, cv.graph)
}

func (cv *Chefviz) main(args []string) {
    recipeRegistory := make(map[string]([]string))
    var recipes []string
    recipes = append(recipes, cv.normalizeRecipeName(args[1]))

    for len(recipes) > 0 {
        recipe := recipes[0]
        _, ok := recipeRegistory[recipe]
        if ok == true {
            // duplicate recipes
            recipes = recipes[1:]
            continue
        }
        recipepath := cv.recipeToFilename(recipe)
        included, err := cv.searchRecipesFromFile(recipepath)
        if err != nil {
            fmt.Println("`" + recipes[0] + "` is not found.")
            recipes = recipes[1:]
            continue
        }

        recipeRegistory[recipe] = included
        recipes = append(recipes, included...)
        cv.addGraph(recipes[0], included)
        recipes = recipes[1:]
    }
    fmt.Println(cv.graph.String())
}

func (cv *Chefviz) normalizeRecipeName (recipe string) string {
    tmp := strings.Split(recipe, "::")
    if len(tmp) > 1 {
        return recipe
    }
    return recipe + "::default"
}

func (cv *Chefviz) recipeToFilename (recipe string) string {
    tmp := strings.Split(recipe, "::")
    return "../sample-chef-repo/cookbooks/" + tmp[0] + "/recipes/" + tmp[1] + ".rb"
}

func (cv *Chefviz) searchRecipesFromFile (filename string) ([]string, error){
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
            ret = append(ret, cv.normalizeRecipeName(recipename[1]))
        }
    }

    return ret, nil
}

func (cv *Chefviz) addGraph(parent string, children []string) {
    cv.graph.AddNode("G", `"`+parent+`"`, nil)
    for _, child := range children {
        cv.graph.AddNode("G", `"`+child+`"`, nil)
        cv.graph.AddEdge(`"`+parent+`"`, `"`+child+`"`, true, nil)
    }
}

