package main

import "fmt"

type VisitorFunc func(*Info, error) error

type Visitor interface {
	Visit(VisitorFunc) error
}

type Info struct {
	Namespace string
	Name string
	OtherThings string
}

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

type OtherThingsVisitor struct {
	visitor Visitor
}

func (v OtherThingsVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("===> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function")
		return err
	})
}

type LogVisitor struct {
	visitor Visitor
}

func (v LogVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("LogVisitor() before call function")
		err = fn(info, err)
		fmt.Println("LogVisitor() after call function")
		return err
	})
}

type NameVisitor struct {
	visitor Visitor
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		err = fn(info, err)
		fmt.Printf("===> name=%s, namespace=%s", info.Name, info.Namespace)
		fmt.Println("NameVisitor() after call function")
		return err
	})
}

// decorate visitor

type DecoratedVisitor struct {
	visitor Visitor
	decorators []VisitorFunc
}

func NewDecoratedVisitor(v Visitor, fn ...VisitorFunc) Visitor {
	if len(fn) == 0 {
		return v
	}

	return DecoratedVisitor{v, fn}
}

func (v DecoratedVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		if err != nil {
			return err
		}
		if err = fn(info, nil); err != nil {
			return err
		}
		for i := range v.decorators {
			if err := v.decorators[i](info, nil); err != nil {
				return err
			}
		}
		return nil
	})
}

func NameVisitor1(info *Info, err error) error {
	fmt.Printf("name=%s, namespace=%s\n", info.Name, info.Namespace)
	
	return nil
}

func OtherVisitor1(info *Info, err error) error {
	fmt.Printf("Other=%s\n", info.OtherThings)

	return nil
}

func main() {
	loadFile := func(info *Info, err error) error {
		info.Name = "fan hu"
		info.Namespace = "wechat"
		info.OtherThings = "we are running as remote team"
		return nil
	}

	//info := Info{}
	//var v Visitor = &info
	//v = LogVisitor{v}
	//v = NameVisitor{v}
	//v = OtherThingsVisitor{v}
	//v.Visit(loadFile)

	// decorate visitor
	info1 := Info{}
	var v1 Visitor = &info1
	v1 = NewDecoratedVisitor(v1, NameVisitor1, OtherVisitor1)
	v1.Visit(loadFile)
}
