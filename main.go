package main

import (
	"coffee-shop/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Base de données en mémoire
var drinks []models.Drink
var orders []models.Order
var orderCounter int = 1

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})

}

func main() {
	drinks = []models.Drink{
		{ID: "D-001", Name: "Espresso", Category: "coffee", BasePrice: 2.0},
		{ID: "D-002", Name: "Cappuccino", Category: "coffee", BasePrice: 3.0},
		{ID: "D-003", Name: "Infusion de Marijuana", Category: "Infusion", BasePrice: 3.5},
		{ID: "D-004", Name: "Latte", Category: "tea", BasePrice: 2.5},
		{ID: "D-005", Name: "Limonade", Category: "cold", BasePrice: 2.0},
		{ID: "D-008", Name: "latte macchiato", Category: "coffee", BasePrice: 3.5},
	}
	// créer le routeur Mux

	// TODO 2 : Créer le routeur mux
	router := mux.NewRouter()

	// Get / Orders

	router.HandleFunc("/orders", getOrders).Methods("GET")
	router.HandleFunc("/menu", getMenu).Methods("GET")

	// GET /orders{id}
	// POST /orders
	router.HandleFunc("/orders", createOrder).Methods("POST")

	// PATCH /orders/{id}/status
	// DELETE /orders/{id}

	// TODO 4 : ( Optionsnel ) Ajouter une routes GET / pour un message de bienveune

	// TODO 5 : Afficher un message indiquant que le serveur démarre
	println("Démarrage du serveur sur le port 8080")

	/// TODO 6 : Démarrer le serveur sur le port 8080
	http.ListenAndServe(":8080", corsMiddleware(router))

}

// GET /menu - Récupérer le menu complet
func getMenu(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type à application/json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drinks)
	// 2. Encoder et retourner le slice drinks en JSON

}

// GET /menu/{id} - Récupérer une boisson spécifique
func getDrink(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")
	// 2. Récupérer l'ID depuis les variables de route (mux.Vars)
	vars := mux.Vars(r)
	id := vars["id"]
	// 3. Parcourir le slice drinks
	for _, drink := range drinks {
		if drink.ID == id {
			// 4. Si trouvé : encoder et retourner la boisson
			json.NewEncoder(w).Encode(drink)
			return
		}
		// 5. Sinon : retourner une erreur 404
		http.Error(w, "erreur 404", http.StatusNotFound)

	}
}

// POST /orders - Créer une nouvelle commande et récupérer la commande créée
// Le formulaire va envoyé du JSON
// Elle envoie directement des données avec le drink id , la taille
func createOrder(w http.ResponseWriter, r *http.Request) { // <--- Id de commande , id boissons etc etc
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")

	// 2. Décoder le body JSON dans une variable Order
	var orderInput models.OrderInput
	var order models.Order

	err := json.NewDecoder(r.Body).Decode(&orderInput)
	if err != nil {
		http.Error(w, "erreur 400", http.StatusBadRequest) // 3. Gérer l'erreur de décodage (400 Bad Request)
		return
	}

	// 4. Vérifier que la boisson (DrinkID) existe dans drinks
	var drink *models.Drink

	for _, d := range drinks {
		if d.ID == orderInput.DrinkID {
			drink = &d
			break
		}
	}
	// 5. Si non trouvée : retourner 400 Bad Request
	if drink == nil {
		http.Error(w, "erreur 400 Bad Request", http.StatusBadRequest)
		return
	}
	// 6. Générer un ID unique (ex: ORD-001, ORD-002...)
	order.ID = fmt.Sprintf("ORD-%03d", orderCounter)
	// 7. Remplir order.DrinkName avec le nom de la boisson
	order.DrinkName = drink.Name
	// 8. Définir order.Status à StatusPending
	order.Status = models.StatusPending
	// 9. Définir order.OrderedAt à time.Now()
	order.OrderedAt = time.Now()
	// 10. Calculer le prix total (appeler calculatePrice)
	order.TotalPrice = calculatePrice(drink.BasePrice, order.Size, order.Extras)
	// 11. Ajouter la commande au slice orders
	orders = append(orders, order)
	// 12. Retourner 201 Created avec la commande en JSON
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)

}

// Fonction helper pour calculer le prix total
func calculatePrice(basePrice float64, size string, extras []string) float64 {

	// TODO:
	// 1. Partir du basePrice
	price := basePrice
	// 2. Ajuster selon la taille:
	//    - "small" : x0.8
	switch size {
	case "small": // Si la taille est de type small alors on multiplie le prix par 0.8
		price *= 0.8
	case "medium":
		price *= 1.0
	case "large":
		price *= 1.3
	}
	// 3. Ajouter 0.50€ pour chaque extra
	price += float64(len(extras)) * 0.5 // On ajoute 0.5 euro par extra avec le type len qui compte le nombre d'extras
	// 4. Retourner le prix total
	return price
}

// GET /orders - Récupérer toutes les commandes
func getOrders(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")
	// 2. Encoder et retourner le slice orders en JSON
	json.NewEncoder(w).Encode(orders)
}

// GET /orders/{id} - Récupérer une commande spécifique
func getOrder(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")
	// 2. Récupérer l'ID depuis les variables de route
	vars := mux.Vars(r)
	// 3. Parcourir le slice orders
	id := vars["id"]
	for _, order := range orders {
		if order.ID == id {
			json.NewEncoder(w).Encode(order)
			return
		}
	}
	// 4. Si trouvé : encoder et retourner la commande
	// 5. Sinon : retourner une erreur 404
	http.Error(w, "erreur 404", http.StatusNotFound)
}

// PATCH /orders/{id}/status - Changer le statut d'une commande
func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Définir le header Content-Type
	w.Header().Set("Content-Type", "application/json")
	// 2. Récupérer l'ID depuis les variables de route
	vars := mux.Vars(r)
	id := vars["id"]
	// 3. Créer une struct temporaire avec un champ Status
	type StatusUpdate struct {
		Status models.OrderStatus `json:"status"`
	}
	var statusUpdate StatusUpdate
	// 4. Décoder le body JSON dans cette struct
	err := json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		http.Error(w, "erreur 400", http.StatusBadRequest) // 5. Gérer l'erreur de décodage
		return
	}

	// 6. Parcourir orders et trouver la commande
	for i, order := range orders { // On utilise l'index i pour pouvoir modifier l'élément dans le slice
		if order.ID == id {
			orders[i].Status = statusUpdate.Status // 7. Mettre à jour le statut de la commande
			json.NewEncoder(w).Encode(orders[i])   // 8. Retourner la commande mise à jour en JSON
			return
		}
	}

	// 9. Si non trouvée : retourner 404
	http.Error(w, "erreur 404", http.StatusNotFound)
}

// DELETE /orders/{id} - Annuler une commande
func deleteOrder(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Récupérer l'ID depuis les variables de route
	// 2. Parcourir orders avec l'index
	// 3. Si la commande est trouvée :
	//    a. Vérifier que le statut n'est pas "picked-up"
	//    b. Si picked-up : retourner 400 Bad Request
	//    c. Sinon : supprimer la commande du slice
	//    d. Retourner 204 No Content
	// 4. Si non trouvée : retourner 404
}
