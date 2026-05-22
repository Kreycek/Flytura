import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { Subject } from 'rxjs';

@Component({
  selector: 'app-modal-msgs',
    
  imports: [CommonModule,TranslateModule],
  templateUrl: './modal-msgs.component.html',
  styleUrl: './modal-msgs.component.css'
})
export class ModalMsgsComponent {

   public isVisible = false;
      public message: string = '';
      public  list:any[] = [];
  
       private responseSubject = new Subject<boolean>();
      // Método para abrir o modal e retornar um Observable
      public openModal(_list:any[],_message:string,_isVisible:boolean): Promise<boolean> {
        
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
      
   fechar() {
      this.responseSubject.next(true);
      this.responseSubject.complete();
      this.isVisible = false;
    }
    
  }
