package banner

import (
	"fmt"
)

// Banner comprises the properties and methods required to render a banner
type Banner struct {
	Name         string
	Description  string
	Version      string
	GoVersion    string
	WebsiteURL   string
	VcsURL       string
	VcsName      string
	EmailAddress string
	BorderChar   string
	lastLen      int
	maxLen       int
}

// WriteSimple is the only exported implementation of a banner
func WriteSimple(b Banner) {
	b.maxLen = b.maxString()
	if b.BorderChar == "" {
		b.BorderChar = "*"
	}

	b.writeNewLine()
	b.writeBar()
	b.writeNameDescription()
	b.writeUnderLine()
	b.writeVersion()
	b.writeGoBuildVersion()
	b.writeWebsite()
	b.writeVCS()
	b.writeEmail()
	b.writeBar()
	b.writeNewLine()
}

func (b Banner) maxString() int {
	mS := 0

	// description
	cS := len(b.Name) + len(" : ") + len(b.Description) + 4
	if cS > mS {
		mS = cS
	}

	// version
	cS = len("Version: ") + len(b.Version) + 4
	if cS > mS {
		mS = cS
	}

	// go build version
	cS = len("Go Build Version: ") + len(b.GoVersion) + 4
	if cS > mS {
		mS = cS
	}

	// website
	cS = len("Website: ") + len(b.WebsiteURL) + 4
	if cS > mS {
		mS = cS
	}

	// vcs
	cS = len(b.VcsName) + len(":  ") + len(b.VcsURL) + 4
	if cS > mS {
		mS = cS
	}

	// email
	cS = len("Email:   ") + len(b.EmailAddress) + 4
	if cS > mS {
		mS = cS
	}

	return mS
}

func (b Banner) writeBar() {
	bar := ""
	for i := 0; i < b.maxLen; i++ {
		bar += b.BorderChar
	}
	fmt.Println(bar)
}

func (b Banner) writeNameDescription() {
	textDesc := fmt.Sprintf("%s - %s", b.Name, b.Description)
	fmt.Println(b.buildLine(textDesc))
}

func (b Banner) writeBlankLine() {
	fmt.Println(b.buildLine(""))
}

func (b Banner) writeNewLine() {
	fmt.Println("")
}

func (b Banner) writeVersion() {
	version := fmt.Sprintf("Version: %s", b.Version)
	fmt.Println(b.buildLine(version))
}

func (b Banner) writeGoBuildVersion() {
	goVersion := fmt.Sprintf("Go Build Version: %s", b.GoVersion)
	fmt.Println(b.buildLine(goVersion))
}

func (b Banner) writeWebsite() {
	website := fmt.Sprintf("Website: %s", b.WebsiteURL)
	fmt.Println(b.buildLine(website))
}

func (b Banner) writeVCS() {
	vcs := fmt.Sprintf("%s:  %s", b.VcsName, b.VcsURL)
	fmt.Println(b.buildLine(vcs))
}

func (b Banner) writeEmail() {
	email := fmt.Sprintf("Email:   %s", b.EmailAddress)
	fmt.Println(b.buildLine(email))
}

func (b Banner) writeUnderLine() {
	underline := ""
	for i := 0; i < b.lastLen; i++ {
		underline += "-"
	}
	fmt.Println(b.buildLine(underline))
}

func (b Banner) buildLine(text string) string {
	pre := fmt.Sprintf("%s ", b.BorderChar)
	post := fmt.Sprintf(" %s", b.BorderChar)
	b.lastLen = len(text)
	padSpace := b.maxLen - len(pre) - len(text) - len(post)
	space := ""
	for i := 0; i < padSpace; i++ {
		space += " "
	}
	line := fmt.Sprintf("%s%s%s%s", pre, text, space, post)

	return line
}
