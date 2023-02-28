package resp

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/2/28
 * Time: 17:26
 * Description: No Description
 */

type Reply interface {
	ToBytes() []byte
}
