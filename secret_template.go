//I used secret.go as my file name, but because they are all of the same package
//it doesn't really matter

package shouldBeMain

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

const (
	BOTTOKEN = "" //Insert Telegram Bot token here
)

func initializeMembers(ctx context.Context) error {
	fmArray := [...]FamilyMember{
		//Insert Family Member Details here
		FamilyMember{ /* Add details here */ },
	}
	for _, fm := range fmArray {
		key := datastore.NewKey(ctx, "FamilyMember", "", int64(fm.Id), nil)
		if _, err := datastore.Put(ctx, key, &fm); err != nil {
			return err
		}
	}

	return nil
}
