package farmer

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"npp_backend/entity/private"
	"npp_backend/entity/public"
	"npp_backend/l10n/translate"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

func (fp *FarmerPsql) GetFavoritePotagers(ctx context.Context, userID int) ([]public.FavoritePotager, *private.Error) {
	rows, err := fp.db.Query(ctx, `
		SELECT
			uf.id, u.username, uf.img_url, uf.commune, ST_AsGeoJSON(uf.coordonnees)::json -> 'coordinates' coordonnees,
			(SELECT COUNT(*) FROM rel_fruits_farmers WHERE farmer_id = fp.farmer_id) fruitsCount,
			(SELECT COUNT(*) FROM rel_graines_farmers WHERE farmer_id = fp.farmer_id) grainesCount,
			(SELECT COUNT(*) FROM rel_legumes_farmers WHERE farmer_id = fp.farmer_id) legumesCount
		FROM
			favorite_potagers fp
		JOIN users u ON TRUE
		JOIN users_farmer uf ON uf.id = fp.farmer_id AND uf.user_id = u.id
		WHERE fp.user_id = $1;
	`, userID)
	if err != nil {
		return nil, &private.Error{
			Location:   "db.farmer.GetFavoritePotagers",
			Line:       29,
			Err:        fmt.Errorf("Potager n°%d introuvable", userID),
			TranslKey:  translate.KeyPotagerNotFound,
			ErrorCode:  1,
			StatusCode: http.StatusNotFound,
		}
	}
	favoritePotagers := []public.FavoritePotager{}
	for rows.Next() {
		favoritePotager := public.FavoritePotager{}
		if err = rows.Scan(&favoritePotager.Farmer.ID,
			&favoritePotager.User.Username,
			&favoritePotager.Farmer.ImgUrl,
			&favoritePotager.Farmer.Commune, &favoritePotager.Farmer.Coordonnees, &favoritePotager.FruitsCount, &favoritePotager.GrainesCount, &favoritePotager.LegumesCount); err != nil {
			return nil, &private.Error{
				Location:   "db.farmer.GetFavoritePotagers",
				Line:       46,
				Err:        fmt.Errorf("Potager n°%d introuvable", userID),
				TranslKey:  translate.KeyInternalServerError,
				ErrorCode:  1,
				StatusCode: http.StatusInternalServerError,
			}
		}
		favoritePotagers = append(favoritePotagers, favoritePotager)
	}
	return favoritePotagers, nil
}

func (fp *FarmerPsql) GetMutedPotagers(ctx context.Context, userID int) ([]public.MutedPotager, *private.Error) {
	rows, err := fp.db.Query(ctx, `
		SELECT
			uf.id, u.username, uf.img_url, uf.commune, ST_AsGeoJSON(uf.coordonnees)::json -> 'coordinates' coordonnees,
			(SELECT COUNT(*) FROM rel_fruits_farmers WHERE farmer_id = mp.farmer_id) fruitsCount,
			(SELECT COUNT(*) FROM rel_graines_farmers WHERE farmer_id = mp.farmer_id) grainesCount,
			(SELECT COUNT(*) FROM rel_legumes_farmers WHERE farmer_id = mp.farmer_id) legumesCount
		FROM
			muted_potagers mp
		JOIN users u ON TRUE
		JOIN users_farmer uf ON uf.id = mp.farmer_id AND uf.user_id = u.id
		WHERE
			mp.user_id = $1;
	`, userID)
	if err != nil {
		return nil, &private.Error{
			Location:   "db.farmer.GetMutedPotagers",
			Line:       74,
			Err:        fmt.Errorf("Potager n°%d introuvable", userID),
			TranslKey:  translate.KeyPotagerNotFound,
			ErrorCode:  1,
			StatusCode: http.StatusNotFound,
		}
	}
	mutedPotagers := []public.MutedPotager{}
	for rows.Next() {
		mutedPotager := public.MutedPotager{}
		if err = rows.Scan(&mutedPotager.Farmer.ID, &mutedPotager.User.Username, &mutedPotager.Farmer.ImgUrl, &mutedPotager.Farmer.Commune, &mutedPotager.Farmer.Coordonnees, &mutedPotager.FruitsCount, &mutedPotager.GrainesCount, &mutedPotager.LegumesCount); err != nil {
			return nil, &private.Error{
				Location:   "db.farmer.GetMutedPotagers",
				Line:       87,
				Err:        fmt.Errorf("Potager n°%d introuvable", userID),
				TranslKey:  translate.KeyPotagerNotFound,
				ErrorCode:  1,
				StatusCode: http.StatusNotFound,
			}
		}
		mutedPotagers = append(mutedPotagers, mutedPotager)
	}
	return mutedPotagers, nil
}

func (fp *FarmerPsql) GetPotager(ctx context.Context, userID int, farmerID int) (*public.Potager, *private.Error) {
	potager := public.Potager{}
	description := &sql.NullString{}
	err := fp.db.QueryRow(ctx, `
		SELECT uf.id, u.username, uf.img_url, uf.description, uf.commune, ST_AsGeoJSON(uf.coordonnees)::json->'coordinates' coordonnees,
		(SELECT EXISTS(SELECT * FROM favorite_potagers WHERE user_id = $2 AND farmer_id = uf.id)) favorite
		FROM users u
		JOIN users_farmer uf ON uf.user_id = u.id AND uf.temporary_disabled = FALSE
		WHERE
		uf.id = $1
		AND
		(SELECT NOT EXISTS(SELECT * FROM muted_potagers WHERE user_id = $2 AND farmer_id = $1))
	`, farmerID, userID).Scan(&potager.Farmer.ID, &potager.User.Username, &potager.Farmer.ImgUrl, description, &potager.Farmer.Commune, &potager.Farmer.Coordonnees, &potager.Farmer.Favorite)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &private.Error{
				Location:   "db.farmer.GetPotager",
				Line:       116,
				Err:        fmt.Errorf("Potager n°%d introuvable", farmerID),
				TranslKey:  translate.KeyPotagerNotFound,
				ErrorCode:  1,
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &private.Error{
			Location:   "db.farmer.GetPotager",
			Line:       125,
			Err:        fmt.Errorf("%w: potager n°%d", err, farmerID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	if description.Valid {
		potager.Farmer.Description = description.String
	}
	// FRUITS
	rows, err := fp.db.Query(ctx, `
		SELECT f.id, f.img_url, f.nom, f.variete, f.systeme_echange, f.prix, f.unite_mesure, f.stock
		FROM fruits f
		JOIN rel_fruits_farmers rff ON rff.fruit_id = f.id
		JOIN users_farmer uf ON uf.id = rff.farmer_id AND uf.temporary_disabled = FALSE
		WHERE uf.id = $1
		AND
		(SELECT NOT EXISTS(SELECT * FROM muted_potagers WHERE user_id = $2 AND farmer_id = $1))
	`, farmerID, userID)
	if err != nil {
		return nil, &private.Error{
			Location:   "db.farmer.GetPotager",
			Line:       149,
			Err:        fmt.Errorf("%w: potager n°%d", err, farmerID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	fruits := []public.Fruit{}
	for rows.Next() {
		fruit := public.Fruit{}
		prix := &sql.NullFloat64{}
		uniteMesure := &sql.NullInt64{}
		systemEchange := &pgtype.Int2Array{}
		if err = rows.Scan(&fruit.ID, &fruit.ImgUrl, &fruit.Nom, &fruit.Variete, systemEchange, prix, uniteMesure, &fruit.Stock); err != nil {
			fmt.Println(err)
			return nil, &private.Error{
				Location:   "db.farmer.GetPotager",
				Line:       165,
				Err:        fmt.Errorf("%w: potager n°%d", err, farmerID),
				TranslKey:  translate.KeyInternalServerError,
				ErrorCode:  1,
				StatusCode: http.StatusInternalServerError,
			}
		}
		for _, value := range systemEchange.Elements {
			fruit.SystemeEchange = append(fruit.SystemeEchange, int(value.Int))
		}

		if prix.Valid {
			fruit.Prix = prix.Float64
		}
		if uniteMesure.Valid {
			fruit.UniteMesure = int(uniteMesure.Int64)
		}

		fruits = append(fruits, fruit)
	}
	potager.Fruits = fruits
	// GRAINES
	rows, err = fp.db.Query(ctx, `
		SELECT g.id, g.img_url, g.nom, g.variete, g.systeme_echange, g.prix, g.stock
		FROM graines g
		JOIN rel_graines_farmers rgf ON rgf.graine_id = g.id
		JOIN users_farmer uf ON uf.id = rgf.farmer_id AND uf.temporary_disabled = FALSE
		WHERE uf.id = $1
		AND
		(SELECT NOT EXISTS(SELECT * FROM muted_potagers WHERE user_id = $2 AND farmer_id = $1))
	`, farmerID, userID)
	if err != nil {
		return nil, &private.Error{
			Location:   "db.farmer.GetPotager",
			Line:       195,
			Err:        fmt.Errorf("%w: potager n°%d", err, farmerID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	graines := []public.Graine{}
	for rows.Next() {
		graine := public.Graine{}
		systemEchange := &pgtype.Int2Array{}
		prix := &sql.NullFloat64{}
		if err = rows.Scan(&graine.ID, &graine.ImgUrl, &graine.Nom, &graine.Variete, systemEchange, prix, &graine.Stock); err != nil {
			return nil, &private.Error{
				Location:   "db.farmer.GetPotager",
				Line:       210,
				Err:        fmt.Errorf("%w: potager n°%d", err, farmerID),
				TranslKey:  translate.KeyInternalServerError,
				ErrorCode:  1,
				StatusCode: http.StatusInternalServerError,
			}
		}
		for _, value := range systemEchange.Elements {
			graine.SystemeEchange = append(graine.SystemeEchange, int(value.Int))
		}

		if prix.Valid {
			graine.Prix = prix.Float64
		}

		graines = append(graines, graine)
	}
	potager.Graines = graines
	// LEGUMES
	rows, err = fp.db.Query(ctx, `
		SELECT l.id, l.img_url, l.nom, l.variete, l.systeme_echange, l.prix, l.unite_mesure, l.stock
		FROM legumes l
		JOIN rel_legumes_farmers rlf ON rlf.legume_id = l.id
		JOIN users_farmer uf ON uf.id = rlf.farmer_id AND uf.temporary_disabled = FALSE
		WHERE uf.id = $1
		AND
		(SELECT NOT EXISTS(SELECT * FROM muted_potagers WHERE user_id = $2 AND farmer_id = $1))
	`, farmerID, userID)
	if err != nil {
		return nil, &private.Error{
			Location:   "db.farmer.GetPotager",
			Line:       236,
			Err:        fmt.Errorf("%w: potager n°%d", err, farmerID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	legumes := []public.Legume{}
	for rows.Next() {
		legume := public.Legume{}
		systemEchange := &pgtype.Int2Array{}
		prix := &sql.NullFloat64{}
		uniteMesure := &sql.NullInt64{}
		if err = rows.Scan(&legume.ID, &legume.ImgUrl, &legume.Nom, &legume.Variete, systemEchange, prix, uniteMesure, &legume.Stock); err != nil {
			return nil, &private.Error{
				Location:   "db.farmer.GetPotager",
				Line:       252,
				Err:        fmt.Errorf("%w: potager n°%d", err, farmerID),
				TranslKey:  translate.KeyInternalServerError,
				ErrorCode:  1,
				StatusCode: http.StatusInternalServerError,
			}
		}

		for _, value := range systemEchange.Elements {
			legume.SystemeEchange = append(legume.SystemeEchange, int(value.Int))
		}

		if prix.Valid {
			legume.Prix = prix.Float64
		}
		if uniteMesure.Valid {
			legume.UniteMesure = int(uniteMesure.Int64)
		}

		legumes = append(legumes, legume)
	}
	potager.Legumes = legumes
	rows.Close()
	return &potager, nil
}

func (fp *FarmerPsql) GetNearestAliments(ctx context.Context, userID int, userCoord []float64, search string) ([]public.NearestAliment, *private.Error) {
	rows, err := fp.db.Query(ctx, `
		WITH potager_area AS (
		SELECT id FROM users_farmer uf
		WHERE
			uf.temporary_disabled = FALSE AND user_id <> $1
			AND
			ST_DWithin(uf.coordonnees::geography, ST_SetSRID(ST_Point($2, $3), 4326)::geography, 15000)
			AND
			(SELECT NOT EXISTS(SELECT * FROM muted_potagers WHERE user_id = $1 AND farmer_id = uf.id))
		), fruits_list AS (
			SELECT
				f.id, f.img_url, f.nom, f.variete, f.systeme_echange, f.prix, f.unite_mesure, f.stock,
				uf.id AS farmer_id, u.username, uf.commune, ST_AsGeoJSON(uf.coordonnees)::json->'coordinates' coordonnees
			FROM
				fruits f
				JOIN users u ON TRUE
				JOIN users_farmer uf ON uf.user_id = u.id
				JOIN rel_fruits_farmers rff ON rff.fruit_id = f.id AND rff.farmer_id = uf.id
			WHERE
				rff.farmer_id IN(SELECT * FROM potager_area) AND SIMILARITY (f.nom, $4) >.3
		), legumes_list AS (
			SELECT
				l.id, l.img_url, l.nom, l.variete, l.systeme_echange, l.prix, l.unite_mesure, l.stock,
				uf.id AS farmer_id, u.username, uf.commune, ST_AsGeoJSON(uf.coordonnees)::json->'coordinates' coordonnees
			FROM
				legumes l
				JOIN users u ON TRUE
				JOIN users_farmer uf ON uf.user_id = u.id
				JOIN rel_legumes_farmers rlf ON rlf.legume_id = l.id AND rlf.farmer_id = uf.id
			WHERE
				rlf.farmer_id IN(SELECT* FROM potager_area) AND SIMILARITY (l.nom, $4) >.3
		), graines_list AS (
			SELECT
				g.id, g.img_url, g.nom, g.variete, g.systeme_echange, g.prix, g.stock,
				uf.id AS farmer_id, u.username, uf.commune, ST_AsGeoJSON(uf.coordonnees)::json->'coordinates' coordonnees
			FROM
				graines g
				JOIN users u ON TRUE
				JOIN users_farmer uf ON uf.user_id = u.id
				JOIN rel_graines_farmers rgf ON rgf.graine_id = g.id AND rgf.farmer_id = uf.id
			WHERE
				rgf.farmer_id IN(SELECT * FROM potager_area) AND SIMILARITY (g.nom, $4) >.3
		)
		SELECT * FROM fruits_list
		UNION ALL
		SELECT * FROM legumes_list
		UNION ALL
		SELECT
			id, img_url, nom, variete, systeme_echange, prix, NULL, stock, farmer_id, username, commune, coordonnees
		FROM
			graines_list
		ORDER BY nom;
	`, userID, userCoord[1], userCoord[0], search)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &private.Error{
				Location:   "db.farmer.GetNearestAliments",
				Line:       331,
				Err:        fmt.Errorf("Aliments les plus proches n°%d introuvable", userID),
				TranslKey:  translate.KeyNearestAlimentsNotFound,
				ErrorCode:  1,
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &private.Error{
			Location:   "db.farmer.GetNearestAliments",
			Line:       340,
			Err:        fmt.Errorf("%w: user n°%d", err, userID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	nearestAliments := []public.NearestAliment{}
	for rows.Next() {
		nearestAliment := public.NearestAliment{}
		prix := &sql.NullFloat64{}
		systemEchange := &pgtype.Int2Array{}
		if err := rows.Scan(&nearestAliment.Aliment.ID, &nearestAliment.Aliment.ImgUrl, &nearestAliment.Aliment.Nom, &nearestAliment.Aliment.Variete, systemEchange, prix, &nearestAliment.Aliment.UniteMesure, &nearestAliment.Aliment.Stock, &nearestAliment.Farmer.ID, &nearestAliment.User.Username, &nearestAliment.Farmer.Commune, &nearestAliment.Farmer.Coordonnees); err != nil {
			return nil, &private.Error{
				Location:   "db.farmer.GetNearestAliments",
				Line:       358,
				Err:        fmt.Errorf("%w: user n°%d", err, userID),
				TranslKey:  translate.KeyInternalServerError,
				ErrorCode:  1,
				StatusCode: http.StatusInternalServerError,
			}
		}
		for _, value := range systemEchange.Elements {
			nearestAliment.Aliment.SystemeEchange = append(nearestAliment.Aliment.SystemeEchange, int(value.Int))
		}
		if prix.Valid {
			nearestAliment.Aliment.Prix = prix.Float64
		}
		nearestAliments = append(nearestAliments, nearestAliment)
	}
	return nearestAliments, nil
}

func (fp *FarmerPsql) GetNearestPotagers(ctx context.Context, userID int, userCoord []float64) ([]public.NearestPotager, *private.Error) {
	rows, err := fp.db.Query(ctx, `
		WITH potager_area AS (
			SELECT id FROM users_farmer uf
			WHERE
				uf.temporary_disabled = FALSE AND user_id <> $1
				AND
				ST_DWithin(uf.coordonnees::geography, ST_SetSRID(ST_Point($2, $3), 4326)::geography, 15000)
				AND
				(SELECT NOT EXISTS(SELECT * FROM muted_potagers WHERE user_id = $1 AND farmer_id = uf.id))
		)
		SELECT
			uf.id, uf.img_url, u.username, uf.commune, ST_AsGeoJSON(uf.coordonnees)::json -> 'coordinates' coordonnees,
			(SELECT EXISTS(SELECT * FROM favorite_potagers WHERE user_id = $1 AND farmer_id = uf.id)) favorite,
			(SELECT COUNT(*) FROM rel_fruits_farmers rff WHERE rff.farmer_id = u.id) fruitsCount,
			(SELECT COUNT(*) FROM rel_graines_farmers rgf WHERE rgf.farmer_id = u.id) grainesCount,
			(SELECT COUNT(*) FROM rel_legumes_farmers rlf WHERE rlf.farmer_id = u.id) legumesCount
		FROM
			users u
		JOIN users_farmer uf ON uf.user_id = u.id
		WHERE
			uf.id IN(SELECT * FROM potager_area)
		ORDER BY
			uf.commune;
	`, userID, userCoord[1], userCoord[0])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &private.Error{
				Location:   "db.farmer.GetNearestPotagers",
				Line:       408,
				Err:        fmt.Errorf("Potagers les plus proches n°%d introuvable", userID),
				TranslKey:  translate.KeyNearestPotagersNotFound,
				ErrorCode:  1,
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &private.Error{
			Location:   "db.farmer.GetNearestPotagers",
			Line:       417,
			Err:        fmt.Errorf("%w: user n°%d", err, userID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	nearestPotagers := []public.NearestPotager{}
	for rows.Next() {
		potager := public.NearestPotager{}

		if err := rows.Scan(
			&potager.Farmer.ID, &potager.Farmer.ImgUrl, &potager.User.Username,
			&potager.Farmer.Commune, &potager.Farmer.Coordonnees, &potager.Farmer.Favorite,
			&potager.FruitsCount, &potager.GrainesCount, &potager.LegumesCount); err != nil {
			return nil, &private.Error{
				Location:   "db.farmer.GetNearestPotagers",
				Line:       431,
				Err:        fmt.Errorf("%w: user n°%d", err, userID),
				TranslKey:  translate.KeyInternalServerError,
				ErrorCode:  1,
				StatusCode: http.StatusInternalServerError,
			}
		}

		nearestPotagers = append(nearestPotagers, potager)
	}
	rows.Close()
	return nearestPotagers, nil
}
