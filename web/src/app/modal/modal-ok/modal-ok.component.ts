import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { Subject } from 'rxjs';

@Component({
  selector: 'app-modal-ok',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './modal-ok.component.html',
  styleUrl: './modal-ok.component.css'
})
export class ModalOkComponent {

    isVisible = false;
    message: SafeHtml = '';
  
    private responseSubject = new Subject<boolean>();
  // Método para abrir o modal e retornar um Observable
  
  constructor(private sanitizer: DomSanitizer) {}

  openModal(message: string, isVisible: boolean): Promise<boolean> {
    this.isVisible = isVisible;
    this.message = this.sanitizer.bypassSecurityTrustHtml(message);
    return new Promise(resolve => {
      this.responseSubject = new Subject<boolean>();
      this.responseSubject.subscribe(response => {
        this.isVisible = false;
        resolve(response);
      });
    });
  }


  confirm() {
    this.isVisible = false;
    this.responseSubject.next(true);
    this.responseSubject.complete();
  }
  
}
