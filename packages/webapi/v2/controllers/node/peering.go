package node

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/webapi/v2/apierrors"
	"github.com/iotaledger/wasp/packages/webapi/v2/models"
)

func (c *Controller) getRegisteredPeers(e echo.Context) error {
	peers := c.peeringService.GetRegisteredPeers()
	peerModels := make([]models.PeeringNodeStatusResponse, len(peers))

	for k, v := range peers {
		peerModels[k] = models.PeeringNodeStatusResponse{
			IsAlive:   v.IsAlive,
			NetID:     v.NetID,
			NumUsers:  v.NumUsers,
			PublicKey: v.PublicKey.String(),
			IsTrusted: v.IsTrusted,
		}
	}

	return e.JSON(http.StatusOK, peerModels)
}

func (c *Controller) getTrustedPeers(e echo.Context) error {
	peers, err := c.peeringService.GetTrustedPeers()
	if err != nil {
		return apierrors.InternalServerError(err)
	}

	peerModels := make([]models.PeeringNodeIdentityResponse, len(peers))
	for k, v := range peers {
		peerModels[k] = models.PeeringNodeIdentityResponse{
			NetID:     v.NetID,
			PublicKey: v.PublicKey.String(),
			IsTrusted: v.IsTrusted,
		}
	}

	return e.JSON(http.StatusOK, peerModels)
}

func (c *Controller) getIdentity(e echo.Context) error {
	self := c.peeringService.GetIdentity()

	peerModel := models.PeeringNodeIdentityResponse{
		NetID:     self.NetID,
		PublicKey: self.PublicKey.String(),
		IsTrusted: self.IsTrusted,
	}

	return e.JSON(http.StatusOK, peerModel)
}

func (c *Controller) checkConnectedPeers(e echo.Context) error {
	var err error
	fmt.Printf("XXX KUKU0 %v\n", time.Now())
	publicKeysStr := strings.Split(e.QueryParam("publicKeys"), ",")
	fmt.Printf("XXX KUKU1 BIND %v\n", publicKeysStr)

	publicKeys := make([]*cryptolib.PublicKey, len(publicKeysStr))
	fmt.Printf("XXX KUKU2 RESP %v\n", publicKeys)
	for i := range publicKeys {
		fmt.Printf("XXX KUKU3 PUBLIC KEY %v\n", i)
		publicKeys[i], err = cryptolib.NewPublicKeyFromString(publicKeysStr[i])
		fmt.Printf("XXX KUKU4 PUBLIC KEY %v %v %v\n", i, publicKeys[i], err)
		if err != nil {
			return apierrors.InvalidPropertyError("publicKey", err)
		}
	}
	fmt.Printf("XXX KUKU5 PUBLIC KEYS %v\n", publicKeys)

	connections := c.peeringService.CheckConnectedPeers(e.Request().Context(), publicKeys)
	fmt.Printf("XXX KUKU6 CONNECTIONS %v\n", connections)

	i := 0
	connectedPeersModel := models.PeeringConnectedResponse{
		Sources: make([]models.PeeringConnectedResponseSinglePeer, len(connections)),
	}
	fmt.Printf("XXX KUKU7 RESP INIT %v\n", connectedPeersModel)
	for publicKeyKey, destinations := range connections {
		connectedPeersModel.Sources[i] = models.PeeringConnectedResponseSinglePeer{
			PublicKey:    publicKeyKey.AsPublicKey().String(),
			Destinations: make([]models.PeeringConnectedResponseSinglePeerDestination, len(destinations)),
		}
		j := 0
		for dPublicKeyKey, err := range destinations {
			connectedPeersModel.Sources[i].Destinations[j] = models.PeeringConnectedResponseSinglePeerDestination{
				PublicKey: dPublicKeyKey.AsPublicKey().String(),
			}
			if err != nil {
				connectedPeersModel.Sources[i].Destinations[j].FailReason = err.Error()
			}
			j++
		}
		i++
	}
	fmt.Printf("XXX KUKU7 RESP FINAL %v\n", connectedPeersModel)

	return e.JSON(http.StatusOK, connectedPeersModel)
}

func (c *Controller) trustPeer(e echo.Context) error {
	var trustedPeer models.PeeringTrustRequest

	if err := e.Bind(&trustedPeer); err != nil {
		return apierrors.InvalidPropertyError("body", err)
	}

	publicKey, err := cryptolib.NewPublicKeyFromString(trustedPeer.PublicKey)
	if err != nil {
		return apierrors.InvalidPropertyError("publicKey", err)
	}

	_, err = c.peeringService.TrustPeer(publicKey, trustedPeer.NetID)
	if err != nil {
		return apierrors.InternalServerError(err)
	}

	return e.NoContent(http.StatusOK)
}

func (c *Controller) distrustPeer(e echo.Context) error {
	var trustedPeer models.PeeringTrustRequest

	if err := e.Bind(&trustedPeer); err != nil {
		return apierrors.InvalidPropertyError("body", err)
	}

	publicKey, err := cryptolib.NewPublicKeyFromString(trustedPeer.PublicKey)
	if err != nil {
		return apierrors.InvalidPropertyError("publicKey", err)
	}

	_, err = c.peeringService.DistrustPeer(publicKey)
	if err != nil {
		return apierrors.InternalServerError(err)
	}

	return e.NoContent(http.StatusOK)
}
