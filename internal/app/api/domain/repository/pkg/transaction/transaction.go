//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../../../test/mock/domain/repository/pkg/$GOPACKAGE/$GOFILE
package transaction

import "context"

type TransactionObject interface {
	Transaction(context.Context, func(context.Context) error) error
}
