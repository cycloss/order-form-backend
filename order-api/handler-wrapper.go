package main

import (
	"github.com/cycloss/aj-bell-test/share"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handlerWrapper struct {
	db *gorm.DB
}

func newHandleWrapper(db *gorm.DB) *handlerWrapper {
	return &handlerWrapper{db: db}
}

type requestProcessorConstructor func(*gin.Context, *gorm.DB, int) requestProcessor

type requestProcessor interface {
	process() (any, error)
}

// Creates and passes a gorm transaction and gin context to the `requestProcessor` created by the `requestProcessorFactory`.
// If the `requestProcessorFactory`'s `process` method returns an error, the transaction is rolled back and
// an response will be sent to the client based on the error.
// If the `requestProcessorFactory`'s `process` method returns no error, and the data it returns is not nil,
// this will be serialised into JSON and returned to the client
func (hw *handlerWrapper) transact(c *gin.Context, rpc requestProcessorConstructor) {
	hw.db.Transaction(func(tx *gorm.DB) error {
		res, err := func() (any, error) {
			claims, err := share.GetUnverifiedJwtClaimsFromHeader(c)
			if err != nil {
				return nil, err
			}
			requestProcessor := rpc(c, tx, claims.ClientId)
			return requestProcessor.process()
		}()
		share.RespondJson(c, res, err)
		return err
	})
}
