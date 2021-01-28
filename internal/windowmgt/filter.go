package windowmgt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aacebedo/snapi3/internal"
	"github.com/rotisserie/eris"
)

type Filter interface {
	IsMatching(element interface{}) bool
	Operator() internal.FilterOperator
}

type WindowFilter struct {
	regex    *regexp.Regexp
	winProp  internal.WindowProperty
	operator internal.FilterOperator
}

func NewWindowFilterWithStr(filterStr string) (res *WindowFilter, err error) {
	properties := []string{string(internal.WindowID), string(internal.WindowClass), string(internal.WindowName), string(internal.WindowType)}
	filterStrRegexStr := fmt.Sprintf("^(?P<winprop>%s):(?P<regex>.*)$", strings.Join(properties, "|"))

	filterStrRegex, err := regexp.Compile(filterStrRegexStr)
	if err != nil {
		err = eris.Wrapf(internal.InvalidArgumentError, "Impossible to compile regex '%s'", filterStrRegexStr)

		return
	}

	filterElements := filterStrRegex.FindStringSubmatch(filterStr)

	if len(filterElements) != 3 {
		err = eris.Wrapf(internal.InvalidArgumentError, "Invalid filter definition '%s'", filterStr)

		return
	}

	res, err = NewWindowFilter(internal.WindowProperty(filterElements[1]), filterElements[2], internal.Or)

	return
}

func NewWindowFilter(winProp internal.WindowProperty, regexStr string,
	operator internal.FilterOperator) (res *WindowFilter, err error) {
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		err = eris.Wrapf(internal.InvalidArgumentError, "Invalid Regex '%s'", regexStr)

		return
	}

	res = &WindowFilter{regex: regex, winProp: winProp, operator: operator}

	return
}

func (sf *WindowFilter) Operator() (res internal.FilterOperator) {
	return sf.operator
}

func (sf *WindowFilter) IsMatching(win *Window) (res bool) {
	res = false

	switch sf.winProp {
	case internal.WindowID:
		res = sf.regex.MatchString(fmt.Sprintf("%#x", win.XWinID()))
	case internal.WindowClass:
		res = sf.regex.MatchString(win.Class())
	case internal.WindowName:
		res = sf.regex.MatchString(win.Name())
	case internal.WindowType:
		res = false

		winTypes := win.Types()
		for _, winType := range winTypes {
			res = res || sf.regex.MatchString(winType)
		}
	}

	return
}

func (sf *WindowFilter) ConvertToConfig() (res internal.FilterConfiguration) {
	res.WinProperty = sf.winProp
	res.Regex = sf.regex.String()
	res.Operator = sf.Operator()

	return
}

type GroupFilter struct {
	groupProp internal.GroupProperty
	regex     *regexp.Regexp
	operator  internal.FilterOperator
}

func NewGroupFilterWithStr(filterStr string) (res *GroupFilter, err error) {
	properties := []string{string(internal.GroupID), string(internal.GroupName)}
	filterStrRegexStr := fmt.Sprintf("^(?P<groupprop>%s):(?P<regex>.*)$", strings.Join(properties, "|"))

	filterStrRegex, err := regexp.Compile(filterStrRegexStr)
	if err != nil {
		err = eris.Wrapf(internal.InvalidArgumentError, "Impossible to compile regex '%s'", filterStrRegexStr)

		return
	}

	filterElements := filterStrRegex.FindStringSubmatch(filterStr)

	if len(filterElements) != 3 {
		err = eris.Wrapf(internal.InvalidArgumentError, "Invalid filter definition '%s'", filterStr)

		return
	}

	res, err = NewGroupFilter(internal.GroupProperty(filterElements[1]), filterElements[2], internal.Or)

	return
}

func NewGroupFilter(groupProp internal.GroupProperty, regexStr string,
	operator internal.FilterOperator) (res *GroupFilter, err error) {
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		err = eris.Wrapf(internal.InvalidArgumentError, "Invalid Regex '%s'", regexStr)

		return
	}

	res = &GroupFilter{regex: regex, operator: operator, groupProp: groupProp}

	return
}

func (sf *GroupFilter) Operator() (res internal.FilterOperator) {
	return sf.operator
}

func (sf *GroupFilter) IsMatching(group *WindowGroup) (res bool) {

	res = false

	switch sf.groupProp {
	case internal.GroupID:
		res = sf.regex.MatchString(fmt.Sprintf("%d", group.ID()))
	case internal.GroupName:
		res = sf.regex.MatchString(group.Name())
	}

	return
}

// type ComposedWindowFilter struct {
// 	filters  *arraylist.List
// 	operator internal.FilterOperator
// }

// func NewComposedWindowFilter(operator internal.FilterOperator) (res *ComposedWindowFilter) {
// 	res = &ComposedWindowFilter{operator: operator, filters: arraylist.New()}

// 	return
// }

// func (cf *ComposedWindowFilter) Operator() (res internal.FilterOperator) {
// 	return cf.operator
// }

// func (cf *ComposedWindowFilter) AddFilter(filterToAdd WindowFilter) {
// 	cf.filters.Add(filterToAdd)
// }

// func (cf *ComposedWindowFilter) IsMatching(win *Window) (res bool) {
// 	res = false

// 	if cf.filters.Size() > 0 {
// 		filterIt := cf.filters.Iterator()
// 		filter := filterIt.Value().(WindowFilter)
// 		res = filter.IsMatching(win)

// 		for filterIt.Next(); filterIt.Next(); {
// 			filter = filterIt.Value().(WindowFilter)
// 			switch filter.Operator() {
// 			case internal.Or:
// 				res = res || filter.IsMatching(win)
// 			case internal.And:
// 				res = res && filter.IsMatching(win)
// 			}
// 		}
// 	}

// 	return
// }

// func (cf *ComposedWindowFilter) ConvertToConfig() (res internal.FilterConfiguration) {
// 	res.Operator = cf.Operator()

// 	for filterIt := cf.filters.Iterator(); filterIt.Next(); {
// 		curFilter := filterIt.Value().(WindowFilter)
// 		filterConfig := curFilter.ConvertToConfig()
// 		res.Filters = append(res.Filters, filterConfig)
// 	}

// 	return
// }
