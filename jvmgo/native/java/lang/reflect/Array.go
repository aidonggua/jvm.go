package reflect

import (
	. "github.com/zxh0/jvm.go/jvmgo/any"
	"github.com/zxh0/jvm.go/jvmgo/jvm/rtda"
	rtc "github.com/zxh0/jvm.go/jvmgo/jvm/rtda/class"
	"github.com/zxh0/jvm.go/jvmgo/native/box"
)

func init() {
	_array(get, "get", "(Ljava/lang/Object;I)Ljava/lang/Object;")
	_array(getLength, "getLength", "(Ljava/lang/Object;)I")
	_array(newArray, "newArray", "(Ljava/lang/Class;I)Ljava/lang/Object;")
}

func _array(method Any, name, desc string) {
	rtc.RegisterNativeMethod("java/lang/reflect/Array", name, desc, method)
}

// public static native Object get(Object array, int index)
//         throws IllegalArgumentException, ArrayIndexOutOfBoundsException;
// (Ljava/lang/Object;I)Ljava/lang/Object;
func get(frame *rtda.Frame) {
	vars := frame.LocalVars()
	arr := vars.GetRef(0)
	index := vars.GetInt(1)

	if arr == nil {
		frame.Thread().ThrowNPE()
		return
	}
	if !arr.IsArray() {
		frame.Thread().ThrowIllegalArgumentException("Argument is not an array")
		return
	}
	if index < 0 || index >= rtc.ArrayLength(arr) {
		frame.Thread().ThrowArrayIndexOutOfBoundsExceptionNoMsg()
		return
	}

	if !arr.IsPrimitiveArray() {
		obj := arr.Refs()[index]
		frame.OperandStack().PushRef(obj)
		return
	}

	// primitive array
	stack := frame.OperandStack()
	primitiveDescriptor := arr.Class().Name()[1]
	switch primitiveDescriptor {
	case 'Z':
		stack.PushBoolean(arr.Booleans()[index] == 1)
	case 'B':
		stack.PushInt(int32(arr.Bytes()[index]))
	case 'C':
		stack.PushInt(int32(arr.Chars()[index]))
	case 'S':
		stack.PushInt(int32(arr.Shorts()[index]))
	case 'I':
		stack.PushInt(arr.Ints()[index])
	case 'J':
		stack.PushLong(arr.Longs()[index])
	case 'F':
		stack.PushFloat(arr.Floats()[index])
	case 'D':
		stack.PushDouble(arr.Doubles()[index])
	}

	// boxing
	box.Box(frame, primitiveDescriptor)
}

// public static native int getLength(Object array) throws IllegalArgumentException;
// (Ljava/lang/Object;)I
func getLength(frame *rtda.Frame) {
	vars := frame.LocalVars()
	arr := vars.GetRef(0)

	// todo IllegalArgumentException
	_len := rtc.ArrayLength(arr)
	stack := frame.OperandStack()
	stack.PushInt(_len)
}

// private static native Object newArray(Class<?> componentType, int length)
// throws NegativeArraySizeException;
// (Ljava/lang/Class;I)Ljava/lang/Object;
func newArray(frame *rtda.Frame) {
	vars := frame.LocalVars()
	componentType := vars.GetRef(0)
	length := vars.GetInt(1)
	if length < 0 {
		// todo
		panic("NegativeArraySizeException")
	}

	componentClass := componentType.Extra().(*rtc.Class)
	arrObj := componentClass.NewArray(uint(length))

	stack := frame.OperandStack()
	stack.PushRef(arrObj)
}
