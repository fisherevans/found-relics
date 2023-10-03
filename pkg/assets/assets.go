package assets

type Resources struct {
	Sprites *Sprites
}

func LoadAll() Resources {
	return Resources{
		Sprites: LoadSprites(),
	}
}
