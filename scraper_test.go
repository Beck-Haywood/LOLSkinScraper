package main

import (
	"testing"
	"reflect"
	"fmt"
)

func TestLinkScrape(t *testing.T){
    got := linkScrape()
    want := Characters{}

    if reflect.DeepEqual(got, want) {
        t.Errorf("got %q, wanted %q", got, want)
    }
}
func TestChampScrape(t *testing.T){
    got := champScrape("https://leagueoflegends.fandom.com/wiki/Rek%27Sai/LoL/Cosmetics")
    want := Character{}

    if reflect.DeepEqual(got, want){
        t.Errorf("got %q, wanted %q", got, want)
    }
}