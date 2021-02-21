package person

import (
	"context"
	"net/http"
	"npp_backend/entity/private"
	"npp_backend/l10n/translate"
	"strings"
)

func (pp *PersonPsql) AddFavoritePotager(ctx context.Context, userId int, farmerId int) *private.Error {
	if _, err := pp.db.Exec(ctx, `
		INSERT INTO favorite_potagers(user_id, farmer_id) VALUES ($1, $2);
	`, userId, farmerId); err != nil {
		if strings.Contains(err.Error(), "favorite_potagers_unique_idx") {
			return &private.Error{
				Location:   "db.person.AddFavoritePotager",
				Line:       16,
				Err:        err,
				TranslKey:  translate.KeyFarmerAlreadyFavorite,
				ErrorCode:  1,
				StatusCode: http.StatusConflict,
			}
		}
		return &private.Error{
			Location:   "db.person.AddFavoritePotager",
			Line:       13,
			Err:        err,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (pp *PersonPsql) AddMutedPotager(ctx context.Context, userId int, farmerId int) *private.Error {
	if _, err := pp.db.Exec(ctx, `
		INSERT INTO muted_potagers(user_id, farmer_id) VALUES ($1, $2);
	`, userId, farmerId); err != nil {
		return &private.Error{
			Location:   "db.person.AddMutedPotager",
			Line:       29,
			Err:        err,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (pp *PersonPsql) DeleteFruit(ctx context.Context, fruitId int) *private.Error {
	if _, err := pp.db.Exec(ctx, `DELETE FROM fruits WHERE id = $1`, fruitId); err != nil {
		return &private.Error{
			Location:   "db.person.DeleteFruit",
			Line:       43,
			Err:        err,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (pp *PersonPsql) DeleteGraine(ctx context.Context, graineId int) *private.Error {
	if _, err := pp.db.Exec(ctx, `DELETE FROM graines WHERE id = $1`, graineId); err != nil {
		return &private.Error{
			Location:   "db.person.DeleteGraine",
			Line:       57,
			Err:        err,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (pp *PersonPsql) DeleteLegume(ctx context.Context, legumeId int) *private.Error {
	if _, err := pp.db.Exec(ctx, `DELETE FROM legumes WHERE id = $1`, legumeId); err != nil {
		return &private.Error{
			Location:   "db.person.DeleteLegume",
			Line:       71,
			Err:        err,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (pp *PersonPsql) DeleteFavoritePotager(ctx context.Context, userId int, farmerId int) *private.Error {
	if _, err := pp.db.Exec(ctx, `
		DELETE FROM favorite_potagers WHERE user_id = $1 AND farmer_id = $2;
	`, userId, farmerId); err != nil {
		return &private.Error{
			Location:   "db.person.DeleteFavoritePotager",
			Line:       87,
			Err:        err,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (pp *PersonPsql) DeleteMutedPotager(ctx context.Context, userId int, farmerId int) *private.Error {
	if _, err := pp.db.Exec(ctx, `
		DELETE FROM muted_potagers WHERE user_id = $1 AND farmer_id = $2
	`, userId, farmerId); err != nil {
		return &private.Error{
			Location:   "db.person.DeleteFavoritePotager",
			Line:       103,
			Err:        err,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}
