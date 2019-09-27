package gq_test

import (
	"github.com/pubgo/g/logs"
	"github.com/pubgo/mycli/internal/gq"
	"testing"
)

func TestParse(t *testing.T) {
	logs.P("parse", gq.Parse(
		`<library>
<!-- Great book. -->
<book id="b0836217462" available="true">
    <isbn>0836217462</isbn>
    <title lang="en">Being a Dog Is a Full-Time Job</title>
    <quote>I'd dog paddle the deepest ocean.</quote>
    <author id="CMS">
        <?echo "go rocks"?>
        <name>Charles M Schulz</name>
        <born>1922-11-26</born>
        <dead>2000-02-12</dead>
    </author>
    <character id="PP">
        <name>Peppermint Patty</name>
        <born>1966-08-22</born>
        <qualification>bold, brash and tomboyish</qualification>
    </character>
    <character id="Snoopy">
        <name>Snoopy</name>
        <born>1950-10-04</born>
        <qualification>extroverted beagle</qualification>
    </character>
</book>
</library>`,
		"{"+
			"bookID `css(\"book\");attr(\"id\")`"+
			"title `css(\"title\")`"+
			"isbn `xpath(\"//isbn\")`"+
			"quote `css(\"quote\")`"+
			"language `css(\"title\");attr(\"lang\")`"+
			"author `css(\"author\")` {"+
			"  name `css(\"name\")`"+
			"    born `css(\"born\")`"+
			"    dead `css(\"dead\")`"+
			"}"+
			"character `xpath(\"//character\")` [{"+
			"    name `css(\"name\")`"+
			"    born `css(\"born\")`"+
			"    qualification `xpath(\"qualification\")`"+
			"}]"+
			"}"))
}
