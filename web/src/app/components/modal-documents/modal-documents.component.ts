import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { Subject } from 'rxjs';

@Component({
  selector: 'app-modal-documents',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './modal-documents.component.html',
  styleUrl: './modal-documents.component.css'
})
export class ModalDocumentsComponent {

  
    isVisible = false;
    message: string = '';
    list:any[] = [];
    
      private responseSubject = new Subject<boolean>();
    // Método para abrir o modal e retornar um Observable
    openModal(_list:any[],_message:string,_isVisible:boolean): Promise<boolean> {
      
      this.list=[];
      this.isVisible = _isVisible;
      this.message=_message;
      this.list=_list;

      return new Promise(resolve => {
        this.responseSubject = new Subject<boolean>();
        this.responseSubject.subscribe(response => {
          this.isVisible = false;
          resolve(response);
        });
      });
    }
  
    confirm() {
      this.responseSubject.next(true);
      this.responseSubject.complete();
      this.isVisible = false;
    }
    
}
