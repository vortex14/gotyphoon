package rod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/go-rod/rod"

	"github.com/vortex14/gotyphoon/interfaces"
)

func AwaitDom(element string) string {
	return fmt.Sprintf(`async () => {

				function waitForElm(selector) {
					return new Promise(resolve => {
						if (document.querySelector(selector)) {
							return resolve(document.querySelector(selector));
						}

						const observer = new MutationObserver(mutations => {
							if (document.querySelector(selector)) {
								resolve(document.querySelector(selector));
								observer.disconnect();
							}
						});

						observer.observe(document.body, {
							childList: true,
							subtree: true
						});
					});
				}


				await waitForElm('%s');

			}`, element)
}

func findCoordsAndClickBy(
	x, y int,
	path string,
	page *rod.Page,
	last bool,
	must bool,
	logger interfaces.LoggerInterface) error {

	err, coords := getCoords(logger, page, path, x, y, last)

	if err != nil {
		logger.Error(err)
		return err
	}

	if must {
		//logger.Debug(coords.Path, " <<")
		page.MustElement(coords.Path).MustClick()
	} else {
		page.Mouse.MustMoveTo(coords.X, coords.Y).MustClick("left")
	}

	return nil
}

type Coords struct {
	X    float64
	Y    float64
	Name string
	Path string
}

type MultiCoords []Coords

type JSTemplate struct {
	X        int
	Y        int
	BasePath string
	DomCssFn string
	IsLast   bool
}

func getDomFullPath() string {
	return `function getDomPath(el) {
		  if (!el) {
			return;
		  }
		  var stack = [];
		  var isShadow = false;
		  while (el.parentNode != null) {
			// console.log(el.nodeName);
			var sibCount = 0;
			var sibIndex = 0;
			// get sibling indexes
			for ( var i = 0; i < el.parentNode.childNodes.length; i++ ) {
			  var sib = el.parentNode.childNodes[i];
			  if ( sib.nodeName == el.nodeName ) {
				if ( sib === el ) {
				  sibIndex = sibCount;
				}
				sibCount++;
			  }
			}
			// if ( el.hasAttribute('id') && el.id != '' ) { no id shortcuts, ids are not unique in shadowDom
			//   stack.unshift(el.nodeName.toLowerCase() + '#' + el.id);
			// } else
			var nodeName = el.nodeName.toLowerCase();
			if (isShadow) {
			  nodeName += "::shadow";
			  isShadow = false;
			}
			if ( sibCount > 1 ) {
			  stack.unshift(nodeName + ':nth-of-type(' + (sibIndex + 1) + ')');
			} else {
			  stack.unshift(nodeName);
			}
			el = el.parentNode;
			if (el.nodeType === 11) { // for shadow dom, we
			  isShadow = true;
			  el = el.host;
			}
		  }
		  stack.splice(0,1); // removes the html element
		  return stack.join(' > ');
		}
`
}

func getCoords(
	logger interfaces.LoggerInterface,
	page *rod.Page,
	cssPath string,
	xPadding int,
	yPadding int,
	isLast bool,
) (error, *Coords) {

	mainTemplate := `() => {

		var getDomFullPath = {{.DomCssFn}}

		if (!document.querySelector('{{.BasePath}}')) { return 0 }
        let element = undefined
		{{if .IsLast}}
			let elements = document.querySelectorAll('{{.BasePath}}')
            let count = elements.length
            element = elements[count-1]
        {{end}}

        {{if not .IsLast}}
			element = document.querySelector('{{.BasePath}}')
		{{end}}

		let X = window.scrollX + element.getBoundingClientRect().left + {{.X}}
		let Y = window.scrollY + element.getBoundingClientRect().top + {{.Y}}
		let Path = getDomFullPath(element)


		return {X,Y,Path}
	}`

	tmpl := template.Must(template.New("new").Parse(mainTemplate))
	templateBuffer := &bytes.Buffer{}
	err := tmpl.Execute(templateBuffer, &JSTemplate{
		BasePath: cssPath,
		X:        xPadding,
		Y:        yPadding,
		DomCssFn: getDomFullPath(),
		IsLast:   isLast,
	})

	source, erp := page.MustEval(templateBuffer.String()).MarshalJSON()

	if erp != nil {
		logger.Error(erp)
		return erp, nil
	}
	coords := &Coords{}
	err = json.Unmarshal(source, coords)
	if err != nil {
		logger.Errorf("%s, %s", err.Error(), string(source))
		return err, nil
	}

	return err, coords
}

func getMultiCoords(
	js string,
	logger interfaces.LoggerInterface,
	page *rod.Page,
	cssPath string,
	xPadding int,
	yPadding int,
	isLast bool,
) (error, MultiCoords) {

	tmpl := template.Must(template.New("new").Parse(js))
	templateBuffer := &bytes.Buffer{}
	err := tmpl.Execute(templateBuffer, &JSTemplate{
		BasePath: cssPath,
		X:        xPadding,
		Y:        yPadding,
		DomCssFn: getDomFullPath(),
		IsLast:   isLast,
	})

	source, erp := page.MustEval(templateBuffer.String()).MarshalJSON()

	if erp != nil {
		logger.Error(erp)
		return erp, nil
	}

	//logger.Debugf("%s", string(source))

	mCoords := make(MultiCoords, 0)
	err = json.Unmarshal(source, &mCoords)
	if err != nil {
		logger.Errorf("%s, %s", err.Error(), string(source))
		return err, nil
	}

	return err, mCoords
}
