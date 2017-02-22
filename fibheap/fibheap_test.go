package fibheap

// Tests for the Fibonacci heap with floating point number priorities

import (
	"testing"

	"math/rand"

	"github.com/stretchr/testify/assert"
)

// Go does not have constant arrays.
// Settling for standard variables.
var NumberSequence1 = [...]float64{6145466173.743959, 1717075442.6908855, -9223106115.008125,
	6664774768.783949, -9185895273.675707, -2271628840.682966, -6843837387.469989,
	-3075112103.982916, -7315786187.596851, 9022422938.330479, 9230482598.051868,
	-2019031911.3141594, 4852342381.928253, 7767018098.497437, -5163143977.984332,
	7265142312.343864, -9974588724.261246, -4721177341.970384, 6608275091.590723,
	-2509051968.8908787, -2608600569.397663, 4602079812.256586, 4204221071.262924,
	2072073006.576254, -1375445006.5510921, 9753983872.378643, 3379810998.918478,
	-2120599284.15699, -9284902029.588614, 3804069225.763077, 4680667479.457649,
	3550845076.5165443, 689351033.7409191, -6170564101.460268, 5769309548.4711685,
	-7203959673.554039, -1542719821.5259266, 8314666872.8992195, 4582459708.761353,
	4558164249.709116, -409019759.7648945, 2050647646.0881348, 3337347280.2468243,
	8841975976.437397, -1540752999.8368673, 4548535015.628077, -7013783667.095476,
	2287926261.9939594, -2539231979.834078, -9359850979.452446, 5390795464.938633,
	-9969381716.563528, 3273172669.620493, -8839719143.511513, 9436856014.244781,
	9032693590.852093, 748366072.01511, -8165322713.346881, -9745450118.0132,
	-6554663739.562494, -8350123090.830288, 4767099194.408716, -741610722.9710865,
	978853190.937952, -4689006449.5764475, 6712607751.828266, 1834187952.9013042,
	8144068220.835762, 2649156704.6132507, 5206492575.513319, 2355676989.886942,
	6014313651.805082, 1559476573.9042358, -611075813.2161636, -3428570708.324188,
	3758297334.844446, -73880069.57582092, 7939090089.227123, -6135368824.336376,
	5680302744.840729, 7067968530.463007, -4736146992.716046, 6787733005.103142,
	8291261997.956814, -7976948033.245457, -2717662205.411746, 1753831326.4953232,
	3313929049.058649, -6798511690.417229, 4259620288.6441, -8795846089.203701,
	666087815.4947224, -3189108786.1266823, 6098522858.07811, 3670419236.2020073,
	-4904172359.7338295, 7081860835.300518, 4838004130.57917, -8403025837.455175,
	2858604246.067789, 9767232443.473625, 1853770486.2323227, 2111315124.8128128,
	-789990089.2266369, 3855299652.837984, -5262051498.344847, 5195097083.198868,
	-9453697711.29756, -144320772.42621613, -3280154832.042288, 4327603656.616592,
	-4916338352.631529, 177342499.89391518, -6863008836.282527, -4462732551.435464,
	563531299.3931465, 243815563.513546, -2177539298.657405, 9064363201.461056,
	7752407089.025448, 5072315736.623476, 1676308335.832735, 2368433225.444128,
	7191228067.770271, -7952866649.176966, 9029961422.270164, -3694580624.20329,
	2396384720.634838, 2919689806.6469193, 2516309466.887434, 5711191379.798178,
	-7111997035.1143055, -5887152915.558975, 7074496594.814234, 72399466.26899147,
	9162739770.93885, 545095642.1330223, 589248875.6552525, 5429718452.359911,
	2670541446.0850983, 7074768275.337322, -9376701618.064901, -719716639.8418808,
	5870465712.600103, 8906050348.824574, 5260686230.481573, 4525930216.3939705,
	-7558925556.569441, -3524217648.1943235, -8559543174.289785, -402353821.38601303,
	-2939238306.2766924, -8421788462.600799, 173509960.46243477, 2823962320.1096497,
	-2040044596.465724, 8093258879.034134, 1026657583.5726833, -5939324535.959578,
	1869187366.0910244, -8488159448.309237, -9162642241.327745, 9198652822.209103,
	9981219597.001732, 1245929264.1492062, 6333145610.418182, -5007933225.524759,
	-7507006648.70326, -8682109235.019928, 7572534048.487186, 9172777289.492256,
	-4374595711.753318, 7302929281.918972, 6813548014.888256, 7839035144.903576,
	-5126801855.122898, 6523728766.098036, -8063474434.226172, -1011764426.4069233,
	-5468146510.412097, -7725685149.169344, 5224407910.623154, 5337833362.662783,
	3878206583.8412895, -9990847539.012056, 2828249626.7454433, -8802730816.790993,
	-6223950138.847174, -5003095866.683969, 3701841328.9391365, -7438103512.551224,
	-1879515137.467103, -6931067459.813007, -3591253518.1452456, -3249229927.5027523,
	249923973.47061348, -7291235820.978601, -4073015010.864023, -3089932753.657503,
	8220825130.164364}

const Seq1FirstMinimum float64 = -9990847539.012056
const Seq1ThirdMinimum float64 = -9969381716.563528
const Seq1FifthMinimum float64 = -9453697711.29756
const Seq1LastMinimum float64 = 9981219597.001732

var NumberSequence2 = [...]float64{-2901939070.965906, 4539462982.372177, -6222008480.049856,
	-1400427921.5968666, 9866088144.060883, -2943107648.529664, 8985474333.11443,
	9204710651.257133, 5354113876.8447075, 8122228442.770859, -8121418938.303131,
	538431208.3261185, 9913821013.519611, -8722989752.449871, -3091279426.694975,
	7229910558.195713, -2908838839.99403, 2835257231.305996, 3922059795.3656673,
	-9298869735.322557}

const Seq2DecreaseKey1Orig float64 = 9913821013.519611
const Seq2DecreaseKey1Trgt float64 = -8722989752.449871
const Seq2DecreaseKey2Orig float64 = 9866088144.060883
const Seq2DecreaseKey2Trgt float64 = -9698869735.322557
const Seq2DecreaseKey3Orig float64 = 9204710651.257133
const Seq2DecreaseKey3Trgt float64 = -9804710651.257133

var NumberSequence2Sorted = [...]float64{-9804710651.257133, -9698869735.322557, -9298869735.322557,
	-8722989752.449871, -8722989752.449871, -8121418938.303131, -6222008480.049856,
	-3091279426.694975, -2943107648.529664, -2908838839.99403, -2901939070.965906,
	-1400427921.5968666, 538431208.3261185, 2835257231.305996, 3922059795.3656673,
	4539462982.372177, 5354113876.8447075, 7229910558.195713, 8122228442.770859,
	8985474333.11443}

var NumberSequence2Deleted3ElemSorted = [...]float64{-9298869735.322557, -8722989752.449871,
	-8121418938.303131, -6222008480.049856, -3091279426.694975, -2943107648.529664,
	-2908838839.99403, -2901939070.965906, -1400427921.5968666, 538431208.3261185,
	2835257231.305996, 3922059795.3656673, 4539462982.372177, 5354113876.8447075,
	7229910558.195713, 8122228442.770859, 8985474333.11443}

var NumberSequence3 = [...]float64{6015943293.071386, -3878285748.0708866, 8674121166.062424,
	-1528465047.6118088, 7584260716.494843, -373958476.80486107, -6367787695.054295,
	6813992306.719868, 5986097626.907181, 9011134545.052086, 7123644338.268343,
	2646164210.08445, 4407427446.995375, -888196668.2563229, 7973918726.985172,
	-6529216482.09644, 6079069259.51853, -8415952427.784341, -6859960084.757652,
	-502409126.89040375}

var NumberSequence4 = [...]float64{9241165993.258648, -9423768405.578083, 3280085607.6687145,
	-5253703037.682413, 3858507441.2785892, 9896256282.896187, -9439606732.236805,
	3082628799.5320206, 9453124863.59945, 9928066165.458393, 1135071669.4712334,
	6380353457.986282, 8329064041.853199, 2382910730.445751, -8478491750.445316,
	9607469190.690144, 5417691217.440792, -9698248424.421888, -3933774735.280322,
	-5984555343.381466}

var NumberSequenceMerged3And4Sorted = [...]float64{-9698248424.421888, -9439606732.236805,
	-9423768405.578083, -8478491750.445316, -8415952427.784341, -6859960084.757652,
	-6529216482.09644, -6367787695.054295, -5984555343.381466, -5253703037.682413,
	-3933774735.280322, -3878285748.0708866, -1528465047.6118088, -888196668.2563229,
	-502409126.89040375, -373958476.80486107, 1135071669.4712334, 2382910730.445751,
	2646164210.08445, 3082628799.5320206, 3280085607.6687145, 3858507441.2785892,
	4407427446.995375, 5417691217.440792, 5986097626.907181, 6015943293.071386,
	6079069259.51853, 6380353457.986282, 6813992306.719868, 7123644338.268343,
	7584260716.494843, 7973918726.985172, 8329064041.853199, 8674121166.062424,
	9011134545.052086, 9241165993.258648, 9453124863.59945, 9607469190.690144,
	9896256282.896187, 9928066165.458393}

func TestEnqueueDequeueMin(t *testing.T) {
	heap := NewFloatFibHeap()
	for i := 0; i < len(NumberSequence1); i++ {
		heap.Enqueue(NumberSequence1[i])
	}

	var min *Entry
	var err error
	for heap.Size() > 0 {
		min, err = heap.DequeueMin()
		assert.NoError(t, err)
		if heap.Size() == 199 {
			assert.Equal(t, Seq1FirstMinimum, min.Priority)
		}
		if heap.Size() == 197 {
			assert.Equal(t, Seq1ThirdMinimum, min.Priority)
		}
		if heap.Size() == 195 {
			assert.Equal(t, Seq1FifthMinimum, min.Priority)
		}
		if heap.Size() == 0 {
			assert.Equal(t, Seq1LastMinimum, min.Priority)
		}
	}
}

func TestFibHeap_Enqueue_Min(t *testing.T) {
	heap := NewFloatFibHeap()
	for i := 0; i < len(NumberSequence1); i++ {
		heap.Enqueue(NumberSequence1[i])
	}

	min, err := heap.Min()
	assert.NoError(t, err)
	assert.Equal(t, Seq1FirstMinimum, min.Priority)
}

func TestFibHeap_Min_EmptyHeap(t *testing.T) {
	heap := NewFloatFibHeap()

	heap.Enqueue(0)
	heap.DequeueMin()

	// Heap should be empty at this point

	min, err := heap.Min()

	assert.EqualError(t, err, "Trying to get minimum element of empty heap")
	assert.Nil(t, min)
}

func TestFibHeap_DequeueMin_EmptyHeap(t *testing.T) {
	heap := NewFloatFibHeap()
	min, err := heap.DequeueMin()

	assert.EqualError(t, err, "Cannot dequeue minimum of empty heap")
	assert.Nil(t, min)
}

func TestEnqueueDecreaseKey(t *testing.T) {
	heap := NewFloatFibHeap()
	var e1, e2, e3 *Entry
	for i := 0; i < len(NumberSequence2); i++ {
		if NumberSequence2[i] == Seq2DecreaseKey1Orig {
			e1 = heap.Enqueue(NumberSequence2[i])
		} else if NumberSequence2[i] == Seq2DecreaseKey2Orig {
			e2 = heap.Enqueue(NumberSequence2[i])
		} else if NumberSequence2[i] == Seq2DecreaseKey3Orig {
			e3 = heap.Enqueue(NumberSequence2[i])
		} else {
			heap.Enqueue(NumberSequence2[i])
		}
	}

	assert.NotNil(t, e1)
	assert.NotNil(t, e2)
	assert.NotNil(t, e3)

	heap.DecreaseKey(e1, Seq2DecreaseKey1Trgt)
	heap.DecreaseKey(e2, Seq2DecreaseKey2Trgt)
	heap.DecreaseKey(e3, Seq2DecreaseKey3Trgt)

	var min *Entry
	var err error
	for i := 0; i < len(NumberSequence2Sorted); i++ {
		min, err = heap.DequeueMin()
		assert.NoError(t, err)
		assert.Equal(t, NumberSequence2Sorted[i], min.Priority)
	}
}

func TestFibHeap_DecreaseKey_EmptyHeap(t *testing.T) {
	heap := NewFloatFibHeap()

	elem := heap.Enqueue(15)
	heap.DequeueMin()

	// Heap should be empty at this point
	min, err := heap.DecreaseKey(elem, 0)

	assert.EqualError(t, err, "Cannot decrease key in an empty heap")
	assert.Nil(t, min)
}

func TestFibHeap_DecreaseKey_NilNode(t *testing.T) {
	heap := NewFloatFibHeap()
	heap.Enqueue(1)
	min, err := heap.DecreaseKey(nil, 0)

	assert.EqualError(t, err, "Cannot decrease key: given node is nil")
	assert.Nil(t, min)
}

func TestFibHeap_DecreaseKey_LargerNewPriority(t *testing.T) {
	heap := NewFloatFibHeap()
	node := heap.Enqueue(1)
	min, err := heap.DecreaseKey(node, 20)

	assert.EqualError(t, err, "The given new priority: 20, is larger than or equal to the old: 1")
	assert.Nil(t, min)
}

func TestEnqueueDelete(t *testing.T) {
	heap := NewFloatFibHeap()
	var e1, e2, e3 *Entry
	for i := 0; i < len(NumberSequence2); i++ {
		if NumberSequence2[i] == Seq2DecreaseKey1Orig {
			e1 = heap.Enqueue(NumberSequence2[i])
		} else if NumberSequence2[i] == Seq2DecreaseKey2Orig {
			e2 = heap.Enqueue(NumberSequence2[i])
		} else if NumberSequence2[i] == Seq2DecreaseKey3Orig {
			e3 = heap.Enqueue(NumberSequence2[i])
		} else {
			heap.Enqueue(NumberSequence2[i])
		}
	}

	assert.NotNil(t, e1)
	assert.NotNil(t, e2)
	assert.NotNil(t, e3)

	var err error

	err = heap.Delete(e1)
	err = heap.Delete(e2)
	err = heap.Delete(e3)

	var min *Entry
	for i := 0; i < len(NumberSequence2Deleted3ElemSorted); i++ {
		min, err = heap.DequeueMin()
		assert.NoError(t, err)
		assert.Equal(t, NumberSequence2Deleted3ElemSorted[i], min.Priority)
	}
}

func TestFibHeap_Delete_EmptyHeap(t *testing.T) {
	heap := NewFloatFibHeap()

	elem := heap.Enqueue(15)
	heap.DequeueMin()

	// Heap should be empty at this point
	err := heap.Delete(elem)
	assert.EqualError(t, err, "Cannot delete element from an empty heap")
}

func TestFibHeap_Delete_NilNode(t *testing.T) {
	heap := NewFloatFibHeap()
	heap.Enqueue(1)
	err := heap.Delete(nil)
	assert.EqualError(t, err, "Cannot delete node: given node is nil")
}

func TestMerge(t *testing.T) {
	heap1 := NewFloatFibHeap()
	for i := 0; i < len(NumberSequence3); i++ {
		heap1.Enqueue(NumberSequence3[i])
	}

	heap2 := NewFloatFibHeap()
	for i := 0; i < len(NumberSequence4); i++ {
		heap1.Enqueue(NumberSequence4[i])
	}

	heap, err := heap1.Merge(&heap2)
	assert.NoError(t, err)

	var min *Entry
	for i := 0; i < len(NumberSequenceMerged3And4Sorted); i++ {
		min, err = heap.DequeueMin()
		assert.NoError(t, err)
		assert.Equal(t, NumberSequenceMerged3And4Sorted[i], min.Priority)
	}
}

func TestFibHeap_Merge_NilHeap(t *testing.T) {
	var heap FloatingFibonacciHeap
	heap = NewFloatFibHeap()
	newHeap, err := heap.Merge(nil)
	assert.EqualError(t, err, "One of the heaps to merge is nil. Cannot merge")
	assert.Equal(t, newHeap, FloatingFibonacciHeap{})
}

// ***************
// BENCHMARK TESTS
// ***************

/*
Since the e.g. Enqeue operation is constant time,
when go benchmark increases N, the prep time
will increase linearly, but the actual operation
we want to measure will always take the same,
constant amount of time.
This means that on some machines, Go Bench
could try to exponentially increase N in order
to decrease noise in the measurement, but it will
get more and more noise. This can cause a system
to run out of RAM. So be careful if you have a fast PC.
I have removed the b.ResetTimer on constant-time
functions to avoid this negative-feedback loop.
*/

// Runs in O(1) time
func BenchmarkFibHeap_Enqueue(b *testing.B) {

	heap := NewFloatFibHeap()

	slice := make([]float64, 0, b.N)
	for i := 0; i < b.N; i++ {
		slice = append(slice, 2*1E10*(rand.Float64()-0.5))
	}

	for i := 0; i < b.N; i++ {
		heap.Enqueue(slice[i])
	}
}

// Runs in O(log(N)) time
func BenchmarkFibHeap_DequeueMin(b *testing.B) {

	heap := NewFloatFibHeap()

	slice := make([]float64, 0, b.N)
	for i := 0; i < b.N; i++ {
		slice = append(slice, 2*1E10*(rand.Float64()-0.5))
		heap.Enqueue(slice[i])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.DequeueMin()
	}
}

// Runs in O(1) amortized time
func BenchmarkFibHeap_DecreaseKey(b *testing.B) {
	heap := NewFloatFibHeap()

	sliceFlt := make([]float64, 0, b.N)
	sliceE := make([]*Entry, 0, b.N)
	for i := 0; i < b.N; i++ {
		sliceFlt = append(sliceFlt, 2*1E10*(float64(i)-0.5))
		sliceE = append(sliceE, heap.Enqueue(sliceFlt[i]))
	}

	for i := 0; i < b.N; i++ {
		// Shift-decrease keys
		heap.DecreaseKey(sliceE[i], sliceFlt[i]-2E10)
	}
}

// Runs in O(log(N)) time
func BenchmarkFibHeap_Delete(b *testing.B) {
	heap := NewFloatFibHeap()

	sliceFlt := make([]float64, 0, b.N)
	sliceE := make([]*Entry, 0, b.N)
	for i := 0; i < b.N; i++ {
		sliceFlt = append(sliceFlt, 2*1E10*(float64(i)-0.5))
		sliceE = append(sliceE, heap.Enqueue(sliceFlt[i]))
	}

	// Delete runs in log(N) time
	// so safe to reset timer here
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := heap.Delete(sliceE[i])
		assert.NoError(b, err)
	}
}

// Runs in O(1) time
func BenchmarkFibHeap_Merge(b *testing.B) {
	heap1 := NewFloatFibHeap()
	heap2 := NewFloatFibHeap()

	for i := 0; i < b.N; i++ {
		heap1.Enqueue(2 * 1E10 * (rand.Float64() - 0.5))
		heap2.Enqueue(2 * 1E10 * (rand.Float64() - 0.5))
		_, err := heap1.Merge(&heap2)
		assert.NoError(b, err)
	}
}
