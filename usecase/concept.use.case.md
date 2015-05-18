#use case of concept

##0 基本规则

1. 一个变量a可以由concept A来声明，表示变量a的类型必须符合A的约束：
1. 一个变量a可以赋值给concept B声明的变量b，只要a满足B所确立的类型特征约束。
1. 一个对象o在赋值给一个concept A声明的变量a时，并未丢失其类型。concept A仅仅起到限制器的作用，相当于对o的类型做了一个特征切片，在a的作用范围内，只能以A所体现的特征进行操作。o在concept声明的变量间传递过程中，依旧保持其类型，并可以随时随地提取。 

```{code}
concept A {
	...
}
var a A
```

一个泛型定义的变量、形参，只是一个对象的载体，对对象形成一组访问的约束。在这些变量和形参的可见范围内，对访问做出约束。其基本作用是将针对对象的操作确定性地连接起来。这种连接可以是静态的，即在编译期建立，也可以是动态的，即在运行时建立。
这些约束只是对访问的屏蔽，并没有对对象，或者对象的引用的类型产生影响。代码可以在任意时刻获取对象的类型。获取类型可以是静态（编译期）的，也可以是动态（运行时）的。编译器会在编译期设法推断出类型，如果做不到，那么就在运行期提取类型。无论是静态的，还是动态的类型获取，对所编写的代码而言是透明的。

##1 covariance & contravariance

赋值

```{code}
concept A {
	...
}

concept B {
	A
}

type X struct {
	...
}
... // 类型X符合concept B的约束

var a []A
var b []B

x := []X{...}
b = x	// 用x（类型是[]X）赋值给b，X满足B的约束，
		//	而[]X满足[]B的约束，可以赋值
a = b	// b满足a的约束。B满足A的约束，因而[]B满足
		//	[]A的约束（[]B可以用所有操作[]A的方式
		//	操作）。
a = x	// 可以，原因同上
b[0] = a[0] // 不可以，B的约束强于A，a无法保证能够
			//	满足所有b的约束。
b = a	// 不可以，原因同上

// 传递
func f1(p1 A) type(p1) {
	return p1
}

x := X{}
y := f1(x)
var z X = y // 成立。concept A并未改变对象x的类型，
			//	在传递过程中并未丢失类型。

ax := []X{...}
var ay []X = sort(ax) // 同理

```


```{code}
func add(lhd Number, rhd Number) Number {	// #1
	type resType great_type(lhd, rhd)
	var res resType
	lopd, ropd := convert(resType, lhd), convert(resType, rhd)
	return add(lopd + ropd)
}

func add(lhd Number, rhd Number) Number {
	where same_type(lhd, rhd)	// 参数类型约束，函数版本对应参数相同的情况
	return lhd + rhd
}

func add(lhd, rhd Number) Number {	//	等价于上述same_type()版本
	return lhd + rhd
}
```

```{code}
x := []Real{10.5, 7, 5, 0xFF59D237C4, 921873654.22918, Fraction{7,9}}

func sum(v []Real) {
	var s Real = 0.0		// 使用concept定义变量，便于后续根据数组中元素实际类型
							//	扩展。
	// var s fload = 0.0	// 不用具体类型
	for _, val := range v {
		s = add(s, val)		//如上述#1版本的add那样，按实际类型扩展s
	}
}
```
上述代码采用runtime concept。对于sum()而言，在运行时才能确定

```{code}
     func FindAll(a []Equable, v Equable) type(a) {
          r := new(type(a))     // 用a的类型创建一个对象
          for _, d := range a {
               if d.Equal(a) {
                    r = append(r, d)
               }
          }
          return r
     }
```
该函数定义存在问题，形参v的类型不定，无法确定a的元素类型的Equal()能够匹配v的类型。采用什么样的规则可以发现这个问题，并给出错误？