package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi/ecs"
)

func DrawDebug(_ *ecs.ECS, screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.1f\nTPS: %.1f\nVSync: %v", ebiten.ActualFPS(), ebiten.ActualTPS(), ebiten.IsVsyncEnabled()))
}
