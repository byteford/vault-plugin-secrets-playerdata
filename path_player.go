package secretsengine

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type stats struct {
	Strength  int `json:"strength"`
	Dexterity int `json:"dexterity"`
}

type playerDataPlayerEntity struct {
	Class      string `json:"class"`
	Experience int    `json:"experience"`
	Stats      stats  `json:"stats"`
}

func (r *playerDataPlayerEntity) GetLevel() int {
	return int(math.Floor(math.Sqrt(float64(r.Experience))))
}

func (r *playerDataPlayerEntity) toResponceData() map[string]interface{} {
	return map[string]interface{}{
		"class":      r.Class,
		"experience": r.Experience,
		"level":      r.GetLevel(),
		"stats":      "/stats",
	}
}

func pathPlayer(b *playerDataBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: framework.GenericNameRegex("name") + "/stats",
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the playerSec",
					Required:    true,
				},
				"dexterity": {
					Type:        framework.TypeInt,
					Description: "dexterity of the player",
					Required:    true,
				},
				"strength": {
					Type:        framework.TypeInt,
					Description: "strength of the player",
					Required:    true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathPlayerStatsRead,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathPlayerStatsWrite,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathPlayerStatsWrite,
				},
			},
			HelpSynopsis:    pathPlayerStatsHelpSynopsis,
			HelpDescription: pathPlayerStatsHelpDescription,
		},
		{
			Pattern: framework.GenericNameRegex("name") + "/level",
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the playerSec",
					Required:    true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathPlayerReadLevel,
				},
			},
			HelpSynopsis:    pathPlayerStatsHelpSynopsis,
			HelpDescription: pathPlayerStatsHelpDescription,
		},
		{
			Pattern: framework.GenericNameRegex("name") + "/" + framework.GenericNameRegex("key"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the playerSec",
					Required:    true,
				},
				"key": {
					Type:        framework.TypeLowerCaseString,
					Description: "The key to ge info from ",
					Required:    true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathPlayerReadPart,
				},
			},
			HelpSynopsis:    pathPlayerStatsHelpSynopsis,
			HelpDescription: pathPlayerStatsHelpDescription,
		},
		{
			Pattern: framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the playerSec",
					Required:    true,
				},
				"class": {
					Type:        framework.TypeString,
					Description: "class of the player",
					Required:    true,
				},
				"experience": {
					Type:        framework.TypeInt,
					Description: "experience for class",
					Required:    false,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathPlayerRead,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathPlayerWrite,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathPlayerWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathPlayerDelete,
				},
			},
			HelpSynopsis:    pathPlayerHelpSynopsis,
			HelpDescription: pathPlayerHelpDescription,
		},
		{
			Pattern: "?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathPlayerList,
				},
			},
			HelpSynopsis:    pathPlayerListHelpSynopsis,
			HelpDescription: pathPlayerListHelpDescription,
		},
	}
}
func setPlayer(ctx context.Context, s logical.Storage, name string, playerEntity *playerDataPlayerEntity) error {
	entry, err := logical.StorageEntryJSON(name, playerEntity)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("failed to create storage entry for player")
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}
	return nil
}
func (b *playerDataBackend) getPlayer(ctx context.Context, s logical.Storage, name string) (*playerDataPlayerEntity, error) {
	if name == "" {
		return nil, fmt.Errorf("missing player name")
	}

	entry, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var player playerDataPlayerEntity

	if err := entry.DecodeJSON(&player); err != nil {
		return nil, err
	}

	return &player, nil
}

func (b *playerDataBackend) pathPlayerReadPart(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getPlayer(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	field := strings.Title(d.Get("key").(string))

	val := reflect.ValueOf(*entry).FieldByName(field)

	if !val.IsValid() {
		return nil, fmt.Errorf("Unexpected key: %s", d.Get("key"))
	}

	return &logical.Response{
		Data: map[string]interface{}{
			field: fmt.Sprint(val),
		},
	}, nil
}

func (b *playerDataBackend) pathPlayerReadLevel(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getPlayer(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"level": entry.GetLevel(),
		},
	}, nil
}

func (b *playerDataBackend) pathPlayerRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getPlayer(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: entry.toResponceData(),
	}, nil
}

func (b *playerDataBackend) pathPlayerWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("missing player name"), nil
	}

	playerEntry, err := b.getPlayer(ctx, req.Storage, name.(string))
	if err != nil {
		return nil, err
	}

	if playerEntry == nil {
		playerEntry = &playerDataPlayerEntity{}
	}

	createOperation := (req.Operation == logical.CreateOperation)

	if class, ok := d.GetOk("class"); ok {
		playerEntry.Class = class.(string)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing class in role")
	}

	if exp, ok := d.GetOk("experience"); ok {
		playerEntry.Experience = exp.(int)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing experience in role")
	}

	if err := setPlayer(ctx, req.Storage, name.(string), playerEntry); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *playerDataBackend) pathPlayerDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, d.Get("name").(string))
	if err != nil {
		return nil, fmt.Errorf("error deleting playerData role: %w", err)
	}
	return nil, nil
}

func (b *playerDataBackend) pathPlayerList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

const (
	pathPlayerHelpSynopsis    = `Manages the Vault player sec.`
	pathPlayerHelpDescription = `
This path allows you to read and write roles used to generate players.
`
	pathPlayerListHelpSynopsis    = `List the existing player in playerData backend`
	pathPlayerListHelpDescription = `players will be listed.`
)
