package api

import (
	"context"
	"errors"
	"phonebook/internal/types"
	"phonebook/internal/utils"
	phonebookv1 "phonebook/pkg/phonebook/v1"

	"github.com/googleapis/go-type-adapters/adapters"
	"go.etcd.io/bbolt"
)

// Errors
var (
	ErrContactNotFound = errors.New("contact not found")
)

// Store
type Store struct {
	phonebookv1.UnimplementedPhonebookStoreServiceServer
	db *bbolt.DB
}

// NewStore
func NewStore(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0644, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}

	//
	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("contacts"))
		return err
	}); err != nil {
		return nil, err
	}

	//
	return &Store{
		db: db,
	}, nil
}

// GetContact
func (s *Store) GetContact(ctx context.Context, req *phonebookv1.GetContactRequest) (*phonebookv1.GetContactResponse, error) {
	var cc types.Contact
	if err := s.db.View(func(tx *bbolt.Tx) error {
		// get bucket
		bucket := tx.Bucket([]byte("contacts"))

		// get contact
		val := bucket.Get([]byte(req.FullName))
		if val == nil {
			return ErrContactNotFound
		}
		c, err := utils.Decode(val)
		if err != nil {
			return err
		}
		cc = *c
		return nil
	}); err != nil {
		return nil, err
	}

	createdAt, err := adapters.TimeToProtoDateTime(cc.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &phonebookv1.GetContactResponse{
		Contact: &phonebookv1.Contact{
			FullName:  cc.FullName,
			Email:     cc.Email,
			Phone:     cc.Phone,
			CreatedAt: createdAt,
		},
	}, nil
}

// PutContact
func (s *Store) PutContact(ctx context.Context, req *phonebookv1.PutContactRequest) (*phonebookv1.PutContactResponse, error) {
	if err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("contacts"))

		createdAt, err := adapters.ProtoDateTimeToTime(req.Contact.CreatedAt)
		if err != nil {
			return err
		}

		contact := &types.Contact{
			FullName:  req.Contact.FullName,
			Email:     req.Contact.Email,
			Phone:     req.Contact.Phone,
			CreatedAt: createdAt,
		}

		val, err := utils.Encode(contact)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(req.Contact.FullName), val)
	}); err != nil {
		return nil, err
	}
	return &phonebookv1.PutContactResponse{}, nil
}

// DeleteContact
func (s *Store) DeleteContact(ctx context.Context, req *phonebookv1.DeleteContactRequest) (*phonebookv1.DeleteContactResponse, error) {
	if err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("contacts"))

		return bucket.Delete([]byte(req.FullName))
	}); err != nil {
		return nil, err
	}
	return &phonebookv1.DeleteContactResponse{}, nil
}

// ListContacts
func (s *Store) ListContacts(context.Context, *phonebookv1.ListContactsRequest) (*phonebookv1.ListContactsResponse, error) {
	var contacts []*phonebookv1.Contact
	if err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("contacts"))
		return bucket.ForEach(func(k []byte, v []byte) error {
			c, err := utils.Decode(v)
			if err != nil {
				return err
			}

			createdAt, err := adapters.TimeToProtoDateTime(c.CreatedAt)
			if err != nil {
				return err
			}
			contacts = append(contacts, &phonebookv1.Contact{
				FullName:  c.FullName,
				Email:     c.Email,
				Phone:     c.Phone,
				CreatedAt: createdAt,
			})
			return nil
		})
	}); err != nil {
		return nil, err
	}

	return &phonebookv1.ListContactsResponse{
		Contacts: contacts,
	}, nil
}

// Close
func (s *Store) Close() {
	s.db.Close()
}
