import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common'; 
@Component({
  selector: 'app-paginator',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './paginator.component.html',
  styleUrl: './paginator.component.css'
})
export class PaginatorComponent {
  @Input() totalPages: number = 1;
  @Input() currentPage: number = 1;
  @Output() pageChange = new EventEmitter<number>();

  get pagesToShow(): (number | string)[] {
    const pages: (number | string)[] = [];
    const maxPagesToShow = 3; // Quantidade de páginas vizinhas a exibir

    if (this.totalPages <= 7) {
      // Se houver poucas páginas, exibe todas
      for (let i = 1; i <= this.totalPages; i++) {
        pages.push(i);
      }
    } else {
      // Adiciona primeira página sempre
      pages.push(1);

      // Se a página atual for maior que maxPagesToShow + 2, mostra "..."
      if (this.currentPage > maxPagesToShow + 2) {
        pages.push('...');
      }

      // Adiciona páginas próximas à atual
      let start = Math.max(2, this.currentPage - maxPagesToShow);
      let end = Math.min(this.totalPages - 1, this.currentPage + maxPagesToShow);
      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      // Se a última página visível for menor que totalPages - 1, adiciona "..."
      if (this.currentPage < this.totalPages - maxPagesToShow - 1) {
        pages.push('...');
      }

      // Adiciona última página sempre
      pages.push(this.totalPages);
    }
    
    return pages;
  }

  changePage(page: number | string) {
    if (typeof page === 'number' && page >= 1 && page <= this.totalPages) {
      this.pageChange.emit(page);
    }
  }
}
