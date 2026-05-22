import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Subject } from 'rxjs';

@Component({
  selector: 'app-modal-confirmation',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './modal-confirmation.component.html',
  styleUrl: './modal-confirmation.component.css'
})
export class ModalConfirmationComponent {
  isVisible = false;
  message: string = '';
  textBtnSim='Sim';
  textBtnNao='Não'

  private responseSubject = new Subject<boolean>();

  // Método para abrir o modal e retornar um Observable
  openModal(isVisible:boolean, message:string, textBtnSim:string='Sim', textBtnNao:string='Não'): Promise<boolean> {
    this.isVisible = isVisible;
    this.message=message;
    this.textBtnSim=textBtnSim;
    this.textBtnNao=textBtnNao;
    return new Promise(resolve => {
      this.responseSubject = new Subject<boolean>();
      this.responseSubject.subscribe(response => {
        // this.isVisible = false;
        resolve(response);
      });
    });
  }

  confirm() {
    
    this.responseSubject.next(true);
    this.responseSubject.complete();
  }

  cancel() {
    this.responseSubject.next(false);
    this.responseSubject.complete();
    this.isVisible =false;
  }
  
}
