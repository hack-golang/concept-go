/*
#use case of concept

##0 基本规则


0.1 声明concept:
concept_decl : "concept" LNAME "(" type_list ")" "{" concept_body "}";
concept_body : 
	| concept_constrain
	| concept_method
	| concept_function
	| type_decl

concept_constrain : LNAME "(" type_list ")";

concept_method : "func" LNAME "." concept_func;

concept_function : "func" concept_func;

concept_func : LNAME "(" type_list_opt ")" ret_types;

type_list_opt :
	| type_list

ret_types :
	| LNAME
	| "(" type_list ")";

*/

// 示例 #1
	concept Reader(T) {
		func T.Read(p []byte) (n int, err error)
	}

	concept Seeker(T) {
		func T.Seek(offset int64, whence int) (int64, error)
	}

	concept ReadSeeker(T) {
		Reader(T)
		Seeker(T)
	}

	concept comparable(T) {
		func Compare(lhd, rhd T) bool
	}

	concept map(T) {
		container(T)

		type KeyType less
		type ValType
		type IterType iter(KeyType, ValType)

		func T.subscript(KeyType) ValType
		func T.iterator() (IterType, IterType)
	}

	concept convertable(T1,T2) {
		func convert(T1,T2) error
	}

	concept comparable(T1,T2) {
		func compare(T1,T2) int
	}


// 0.2 concept的使用
// 0.2.1 concept约束
// 一个concept约束类型的基本形式是：
//  concept_constrain : LNAME "(" type_list ")";
// concept约束无法单独使用，必须使用在类型声明、函数声明、变量声明等处使用，并且仅用于约束类型。
// 示例：
	type A
	type B
	int(A)				// 约束类型A，A必须满足int所确定的类型
	convertable(B, A)	// 约束类型A和B的关系，两者必须满足convertable，也就是可转换
// 一个concept可以约束任意多个类型。只约束一个类型的称为type concept，使用时可以省略类型参数：
	type A int			// 其含义是：这里需要一个类型A，它满足int所确定的约束。
// 注意：不同于类型的定义：
	type X int32
// X将是一个明确的类型，并且同int32具有相同的结构和方法。但A则不是一个具体的类型，只是一个受到int约束的类型占位符。
//  也就是说，语句 type A intA 不是类型声明，因此不能独立使用，必须在结构、函数的声明体内，以及require子句中使用。
// 可以使用多个concept多次约束一个类型：
	Reader(C)
	Writer(C)
// 表示类型C同时满足Reader和Writer的约束。
// concept可以组合使用，对一个类型做出多重约束：
	Reader(C) && Writer(C)
// 等价于上述多次约束。
// concept嵌套约束：
	type Y Reader(Writer)
// 等价于
	type Y
	Reader(Y) && Writer(Y)
// 嵌套后的concept产生了一个具有更强约束的concept：
	concept ToStringable(T1,T2) convertable(string(T1), T2)
// 等价于
	concept ToStringable(T1,T2) {
		convertable(T1, T2)
		string(T1)
	}
// 只约束单个类型的concept称为type concept，可以省略(T)，直接用concept名
	func Sort(a rand_iterable) {...}	// 形参a必须是一个满足rand_iterable的对象
// 等价于
	func Sort(a T) require rand_iterable(T) {...}

// 0.2.2 结构和函数声明
// 结构声明中可以使用concept约束type：
	type ComplexNumber struct {
		type NumType float

		RealPart	NumType
		ImagePart	NumType
	}

	c1 := ComplexNumber{11.43, 12.7}

	var r1, i1 float64 = 7.62, 5.56
	c2 := ComplexNumber{r1, i1}

// 存在不确定的类型的struct成为泛化的struct，需要在确定类型后才能使用：
	type Complex32 ComplexNumber{ NumType : float32 }
	c3 := Complex32{0.50, 0.30}

// 使用concept约束field：
	type Buffer struct {
		MaxSize		int
		data		[]byte
	}

	b := Buffer{int32(500)}		// MaxSize成为32位整数

// 使用concept约束参数：
	func Max(a fwd_iterable) a.ElemType					// fwd_iterable是concept，意思是这个类型可以迭代（容器）
		require comparable(a.ElemType, a.ElemType) {	// 要求容器元素的类型是可比较的
		...
	}


// 0.2.3 使用concept约束变量：
//	concept可以像类型那样“声明”一个变量：
	var a T
	Reader(T)			// 变量a的类型必须满足Reader的约束
// 使用简写形式
	var a Reader		// 意思是这个变量的类型必须符合concept Reader的要求
	f, _ := os.Open(...)
	var r Reader = f	// r是一个受Reader约束的变量，它具备Reader的特征（接口），只能以
						//	Reader的方式被使用。r被f赋值为，执行的是浅拷贝，仍旧可以控制
						//	f所指代的对象，r的类型仍旧是type(f)。
	var v map = f		// 错误，f不符合concept map的约束

// 一个concept"定义"的变量、形参，只是一个对象的载体，对对象形成一组访问的约束。在这些变量和形参的可见范围内，对访问做出约束。
//	其基本作用是将针对对象的操作确定性地连接起来。这种连接可以是静态的，即在编译期建立，也可以是动态的，即在运行时建立。
//	这些约束只是对访问的屏蔽，并没有对对象，或者对象的引用的类型产生影响。代码可以在任意时刻获取对象的类型。
//	获取类型可以是静态（编译期）的，也可以是动态（运行时）的。编译器会在编译期设法推断出类型，如果做不到，那么就在运行期提取类型。
//	无论是静态的，还是动态的类型获取，对所编写的代码而言是透明的。

// 0.3 类型操作
// 语言层面提供一组类型函数和操作符，用于类型的提取和操作。类型函数和操作有两个实现，分别用于编译期和运行期。
//	相关的变量参数可以在编译期推导出类型的，使用编译期版本，在编译期执行类型操作。否则使用运行期版本，在运行期
//	执行类型操作。运行期的类型操作实际上就是基于反射的操作。
// 0.3.1 type()
// 从一个变量上提取出它的类型。由于不存在继承、子类、多态，一个对象只会有一个类型。
	i := int64(1234567890)
	type(i)						// 返回类型int64
// 可以用于声明新的变量：
	var j type(i) = int64(9876543210)
// 用于函数参数和返回值：
	func Concat(a fwd_iterable, b type(a)) type(a) {
		type ret_t type(a)
		res := ret_t{}

		...						// 合并两个容器

		return res
	}

/*
concept约束组合，操作符： ==, !=, &&, ||, !, ()


1. 一个变量a可以由concept A来声明，表示变量a的类型必须符合A的约束：
1. 一个变量a可以赋值给concept B声明的变量b，只要a满足B所确立的类型特征约束。
1. 一个对象o在赋值给一个concept A声明的变量a时，并未丢失其类型。concept A仅仅起到限制器的作用，相当于对o的类型做了一个特征切片，在a的作用范围内，只能以A所体现的特征进行操作。o在concept声明的变量间传递过程中，依旧保持其类型，并可以随时随地提取。 

concept定义变量未确定类型，runtime也未能落实类型
*/


// 1. 内置concept和函数
//// 任何类型
concept any(T) {}

//// 基础concept
// assignable: 赋值 => a = b
// assignable为语言内置，不可实现，仅用于约束
//	1. 所有assign操作都是数据字节序列的复制，即浅拷贝
//	2. 同一个类型的不同对象可以assignable
//	3. 类型B从类型A直接定义而来，可以将类型B的对象赋值给类型A的对象
//	4. assignable不执行auto bind，由语言本身buildin。（可防止用户实现的assign破坏语义）
concept assignable(T1, T2) {
	T1 == T2 || decl_by(T1, T2) || decl_by(T2, T1)
}

// pointer: 指针 => *<type>
// pointer为语言本身内置，不可实现，仅用于约束
concept pointer(T) {
	type BaseType
	T == *BaseType
}

// copyable: 两个类型间可复制 => copy(a,b)
// copyable也是语言本身内置，不可实现，仅用于约束
// 语义上copyable做+1层的深拷贝，其余规则同assignable
concept copyable(T1, T2) {
	T1 == T2 || decl_by(T1, T2) || decl_by(T2, T1)
}

// convertable: 两个类型的对象间可转换 => x = y.(<type>)
// convertable是单向转换，即：convertable(T1,T2) ≠> convertable(T2,T1)
concept convertable(T1,T2) {
	func convert(T1,T2) error
}

// cloneable: 克隆，执行完全深拷贝
concept cloneable(T) {
	func Clone(T) T
}

// iterator: 迭代器
concept iterator(I) {
	type I.ItemType
	type I.DistType

	func I.value_ref() *I.ItemType		// a := *it; *it = a
}

// fwd_iter: 前向迭代器
concept fwd_iter(I) {
	iterator(I)

	func I.successor() I
}

// bi_iter: 双向迭代器
concept bid_iter(I) {
	fwd_iter(I)

	func I.previous() I
}

// rand_iter: 随机访问迭代器
concept rand_iter(I) {
	bid_iter(I)

	func I.forward(I.DistType) I
	func I.backward(I.DistType) I
	func I.subscpt(I.DistType) *I.ItemType
	func I.subscpt_range(I.DistType, I.DistType) *I.ItemType
	func I.comp(I.DistType) int
}

// iteratable: 可迭代
concept iteratable(T) {
	type T.IterType iterator

	func T.Begin() T.IterType
	func T.End() T.IterType
}

// fwd_iterable: 可前向迭代
concept fwd_iterable(T) {
	iteratable(T)
	fwd_iter(T.IterType)
}

// bid_iterable: 可双向迭代
concept bid_iterable(T) {
	iteratable(T)
	bid_iter(T.IterType)
}

// rand_iterable: 可随机迭代
//	当一个类型满足了随机迭代的约束，便可以以数组的语法使用。
concept rand_iterable(T) {
	iteratable(T)
	rand_iter(T.IterType)
}


//// 数字concept
// number: 数字，所有数字
concept number(T) {

}

// real: 实数
concept real(T) {
	number(T)
}

// rational: 有理数
concept rational(T) {
	real(T)
}

// int: 整型
concept int(T) {
	real(T)
}

// float: 浮点数
concept float(T) {
	real(T)
}

// complex: 复数
concept complex(T) {
	number(T)
}

//// 容器concept
// container: 所有容器
concept container(T) {
	type ElemType
	func len(T) int
	func cap(T) int
	func concat(T, T)
}

// array: 数组。 => [n]xxx，[]xxx, string
concept array(T) {
	container(T)

	type IterType	rand_iter(pair(int, ElemType))

	func T.subscpt(int) ElemType			// => [i]
	func T.subscpt_range(int, int) ElemType	// => [i:j]
	func T.iter() (IterType, IterType)		// => range
}

// list: 链表。
concept list(T) {
	container(T)

	type IterType bi_iter(pair(nil, ElemType))
	func T.iter() (IterType, IterType)		// => range
}

// map: 映射表
concept map(T) {
	container(T)

	type KeyType less
	type ElemType
	type IterType bi_iter(pair(KeyType, ValType))

	func T.subscpt(KeyType) ValType		// => [k]
	func T.iter() (IterType, IterType)	// => range
}


//// 算法
// sort
func Sort(a rand_iterable) {...}
func SortFunc(a rand_iterable, cmp func (a.ElemType, a.ElemType) int) {...}
func Sort(l list(comparable)) {...}
func SortFunc(l list(comparable), cmp func (l.ElemType, l.ElemType) int) {...}

// find
func Find(c fwd_iterable, v c.ElemType) c.ElemType {...}
func FindFunc(c fwd_iterable, e func(c.ElemType, c.ElemType)bool) c.ElemType {...}
func FindAll(c fwd_iterable, v c.ElemType) []c.ElemType {...}
func FindAllFunc(c fwd_iterable, e func(c.ElemType, c.ElemType)bool) []c.ElemType {...}

// split
func Split(c fwd_iterable, v c.ElemType) [][]c.ElemType {...}
func Split(c fwd_iterable, e func(c.ElemType, c.ElemType)bool) [][]c.ElemType {...}

// sum、avg、stddev
func Sum(c fwd_iterable) c.ElemType require addable(c.ElemType) {...}
func Avg(c fwd_iterable) c.ElemType require addable(c.ElemType) {...}
func StdDev(c fwd_iterable) c.ElemType require real(c.ElemType) {...}

// max、min、mid、maxn、minn
func Max(c fwd_iterable) c.ElemType {...}
func MaxFunc(c fwd_iterable, cmp func (c.ElemType, c.ElemType) int) c.ElemType {...}
func Min(c fwd_iterable) c.ElemType require comparable(c.ElemType) {...}
func MinFunc(c fwd_iterable, cmp func (c.ElemType, c.ElemType) int) c.ElemType {...}
func Mid(c fwd_iterable) c.ElemType require comparable(c.ElemType) {...}
func MidFunc(c fwd_iterable, cmp func (c.ElemType, c.ElemType) int) c.ElemType {...}
func MaxN(c fwd_iterable) c.ElemType require comparable(c.ElemType) {...}
func MaxNFunc(c fwd_iterable, cmp func (c.ElemType, c.ElemType) int) c.ElemType {...}
func MinN(c fwd_iterable) c.ElemType require comparable(c.ElemType) {...}
func MinNFunc(c fwd_iterable, cmp func (c.ElemType, c.ElemType) int) c.ElemType {...}



/*
##1 covariance & contravariance

赋值
*/
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


// 对最泛化的数字类型编写加法
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

// 动多态：runtime concept
x := []Real{10.5, 7, 5, 0xFF59D237C4, 921873654.22918, Fraction{7,9}}

func sum(v []Real) {
	var s Real = 0.0		// 使用concept定义变量，便于后续根据数组中元素实际类型
							//	扩展。
	// var s fload = 0.0	// 不用具体类型
	for _, val := range v {
		s = add(s, val)		//如上述#1版本的add那样，按实际类型扩展s
	}
}
/*
上述代码采用runtime concept。对于sum()而言，在运行时才能确定
*/

/*
以下函数定义存在问题，形参v的类型不定，无法确定a的元素类型的Equal()能够匹配v的类型。采用什么样的规则可以发现这个问题，并给出错误？
*/
func FindAll(a []Equable, v Equable) type(a) {
  r := new(type(a))     // 用a的类型创建一个对象
  for _, d := range a {
       if d.Equal(a) {
            r = append(r, d)
       }
  }
  return r
}


/*
深拷贝和浅拷贝：
golang的赋值使用浅拷贝语义。容器的copy采用1层深拷贝语义。
假设实现深拷贝语义：
*/
type A struct {
	f1	B
	f2	string
	f3	[]int
	f4	map[string]int
}

func deepcopy(dst, src A) {
	deepcopy(&dst.f1, &src.f1)
	deepcopy(&dst.f2, &src.f2)
	deepcopy(&dst.f3, &src.f3)
	deepcopy(&dst.f4, &src.f4)
}

type B struct {
	f1	[]string
	f2	map[string]func (float) int
}

func deepcopy(dst, src B) {
	deepcopy(&dst.f1, &src.f1)
	deepcopy(&dst.f2, &src.f2)
}

/* ... */
concept deepcopyable(dstType, srcType) {
	func deepcopy(*dstType, *srcType)
}

func deepcopy(dst, src *array)
	require deepcopyable(dst.ElemType, src.ElemType)) {
	sz_dst, sz_src := len(dst), len(src)
	to_copy := min(sz_dst, sz_src)
	for i := 0; i < to_copy; i++ {
		deepcopy(&dst[i], &src[i])
	}
}

func deepcopy(dst, src *map)
	require deepcopyable(dst.ElemType, src.ElemType)) {
	for key, val := range src {
		var dst_val dst.ElemType
		deepcopy(&dst_val, &val)
		dst[key] = dst_val
	}
}

func deepcopy(dst, src *simple) 		// simple这个concept如何定义？抽象出特征，还是强行bind？
	require assinable(dst, src) {		//	是否提供强行bind这个功能，如何约束这个功能
	*dst = *src
}


/*
强制将一个类型bind到一个concept：
	只能将类型bind到一个没有任何约束的“空”concept。其余类型自动bind到相匹配的concept。
*/

/*
concept assignable：
*/
