package main

import "encoding/xml"

type Manifest struct {
	XMLName     xml.Name    `xml:"manifest"`
	Package     string      `xml:"package,attr"`
	Application Application `xml:"application"`
}

type Application struct {
	XMLName    xml.Name   `xml:"application"`
	Activities []Activity `xml:"activity"`
}

func (a Application) MainActivity() string {
	for _, a := range a.Activities {
		for _, i := range a.IntentFilter {
			for _, n := range i.Actions {
				if n.Name == "android.intent.action.MAIN" {
					return a.Name
				}
			}
		}
	}
	return ""
}

type Activity struct {
	XMLName      xml.Name       `xml:"activity"`
	Name         string         `xml:"name,attr"`
	IntentFilter []IntentFilter `xml:"intent-filter"`
}

type IntentFilter struct {
	XMLName xml.Name `xml:"intent-filter"`
	Actions []Action `xml:"action"`
}

type Action struct {
	XMLName xml.Name `xml:"action"`
	Name    string   `xml:"name,attr"`
}
