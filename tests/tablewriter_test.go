package tests

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"testing"
)


func TestName1(t *testing.T) {
	{
		data := [][]string{
			[]string{"A", "The Good", "500"},
			[]string{"B", "The Very very Bad Man", "288"},
			[]string{"C", "The Ugly", "120"},
			[]string{"D", "The Gopher", "800"},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Sign", "Rating"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
	}

	{
		data := [][]string{
			[]string{"1/1/2014", "Domain name", "2233", "$10.98"},
			[]string{"1/1/2014", "January Hosting", "2233", "$54.95"},
			[]string{"1/4/2014", "February Hosting", "2233", "$51.00"},
			[]string{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
		table.SetFooter([]string{"", "", "Total", "$146.93"}) // Add Footer
		table.SetBorder(false)                                // Set Border to false
		table.AppendBulk(data)                                // Add Bulk Data
		table.Render()
	}

	{
		data := [][]string{
			[]string{"1/1/2014", "Domain name", "2233", "$10.98"},
			[]string{"1/1/2014", "January Hosting", "2233", "$54.95"},
			[]string{"1/4/2014", "February Hosting", "2233", "$51.00"},
			[]string{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
	}

	{
		data := [][]string{
			[]string{"1/1/2014", "Domain name", "1234", "$10.98"},
			[]string{"1/1/2014", "January Hosting", "2345", "$54.95"},
			[]string{"1/4/2014", "February Hosting", "3456", "$51.00"},
			[]string{"1/4/2014", "February Extra Bandwidth", "4567", "$30.00"},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
		table.SetFooter([]string{"", "", "Total", "$146.93"})
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()
	}

	{
		data := [][]string{
			[]string{"1/1/2014", "Domain name", "2233", "$10.98"},
			[]string{"1/1/2014", "January Hosting", "2233", "$54.95"},
			[]string{"1/4/2014", "February Hosting", "2233", "$51.00"},
			[]string{"1/4/2014", "February Extra Bandwidth", "2233", "$30.00"},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "Description", "CV2", "Amount"})
		table.SetFooter([]string{"", "", "Total", "$146.93"}) // Add Footer
		table.SetBorder(false)                                // Set Border to false

		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
			tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold, tablewriter.BgBlackColor},
			tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.FgWhiteColor})

		table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor})

		table.SetFooterColor(tablewriter.Colors{}, tablewriter.Colors{},
			tablewriter.Colors{tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgHiRedColor})

		table.AppendBulk(data)
		table.Render()
	}

	{
		data := [][]string{
			[]string{"A", "The Good", "500"},
			[]string{"B", "The Very very Bad Man", "288"},
			[]string{"C", "The Ugly", "120"},
			[]string{"D", "The Gopher", "800"},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Sign", "Rating"})
		table.SetCaption(true, "Movie ratings.")

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
	}
}
