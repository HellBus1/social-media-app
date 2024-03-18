package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ServerCheck(ginCtx *gin.Context) {
	DB, ok := ginCtx.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	var bank_name string

	err := DB.QueryRow(context.Background(), "select bank_name from bank_accounts limit 1").Scan(&bank_name)
	if err != nil {
		fmt.Printf("QueryRow failed: %v\n", err)
	}

	fmt.Println(bank_name)

	ginCtx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
