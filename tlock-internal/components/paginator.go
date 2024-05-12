package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
)

func Paginator(listview list.Model) string {
	// Total pages
	totalPages := listview.Paginator.TotalPages

	// Paginator items
	paginatorItems := make([]string, totalPages)

	// Add paginator dots
	for index := 0; index < totalPages; index++ {
		renderer := tlockstyles.Styles.SubText.Copy().Bold(true).Render

		if index == listview.Paginator.Page {
			renderer = tlockstyles.Styles.Title.Render
		}

		paginatorItems = append(paginatorItems, renderer("•"))
	}

	// Add to ui
	return lipgloss.JoinHorizontal(lipgloss.Center, paginatorItems...)
}
