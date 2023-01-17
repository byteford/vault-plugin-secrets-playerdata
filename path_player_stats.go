package secretsengine

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (r *stats) toResponceData() map[string]interface{} {
	return map[string]interface{}{
		"Strength":  r.Strength,
		"Dexterity": r.Dexterity,
	}
}

func setPlayerStats(ctx context.Context, s logical.Storage, name string, playerEntity *playerDataPlayerEntity) error {
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
func (b *playerDataBackend) getPlayerStats(ctx context.Context, s logical.Storage, name string) (*playerDataPlayerEntity, error) {
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

func (b *playerDataBackend) pathPlayerStatsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getPlayerStats(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: entry.Stats.toResponceData(),
	}, nil
}

func (b *playerDataBackend) pathPlayerStatsWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("missing player name"), nil
	}

	playerEntry, err := b.getPlayerStats(ctx, req.Storage, name.(string))
	if err != nil {
		return nil, err
	}

	if playerEntry == nil {
		playerEntry = &playerDataPlayerEntity{}
	}

	createOperation := (req.Operation == logical.CreateOperation)

	if strength, ok := d.GetOk("strength"); ok {
		playerEntry.Stats.Strength = strength.(int)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing strength in role")
	}

	if dexterity, ok := d.GetOk("dexterity"); ok {
		playerEntry.Stats.Dexterity = dexterity.(int)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing dexterity in role")
	}

	if err := setPlayerStats(ctx, req.Storage, name.(string), playerEntry); err != nil {
		return nil, err
	}
	return nil, nil
}

const (
	pathPlayerStatsHelpSynopsis    = `Manages the Vault player stats.`
	pathPlayerStatsHelpDescription = `
This path allows you to read and write stats used to generate players.
`
)
