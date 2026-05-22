import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Subject } from 'rxjs';

@Component({
  selector: 'app-models',
    standalone: true,
  imports: [CommonModule],
  templateUrl: './models.component.html',
  styleUrl: './models.component.css'
})
export class ModelsComponent {

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
