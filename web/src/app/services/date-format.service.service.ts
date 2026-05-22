import { Injectable } from '@angular/core';
import { DateAdapter } from '@angular/material/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class DateFormatService {
  private dateFormatSource = new BehaviorSubject(this.getDefaultFormat());
  currentDateFormat = this.dateFormatSource.asObservable();

  constructor(private dateAdapter: DateAdapter<any>) {}

  updateDateFormat(format: any) {
    this.dateAdapter.setLocale('pt'); // Se necessário, ajuste a localidade aqui
    return format;
  }

  // Formato inicial padrão
  private getDefaultFormat() {

    return {
      parse: { dateInput: 'DD/MM/YYYY' },
      display: {
        dateInput: 'DD/MM/YYYY',
        monthYearLabel: 'DD/MMM YYYY',
        dateA11yLabel: 'LL',
        monthYearA11yLabel: 'DD MMMM YYYY'
      }
    };
  }

  // Método para alterar o formato dinamicamente
  changeDateFormat(tp: string) {
    
    if (tp === '2') {

      this.dateFormatSource.next({
        parse: { dateInput: 'DD/MM/YYYY' },
        display: {
          dateInput: 'DD/MM/YYYY',
          monthYearLabel: 'DD/MMM YYYY',
          dateA11yLabel: 'LL',
          monthYearA11yLabel: 'DD MMMM YYYY'
        }
      });
    } else {
  
      this.dateFormatSource.next({
        parse: { dateInput: 'MM/YYYY' },
        display: {
          dateInput: 'MM/YYYY',
          monthYearLabel: 'MMM YYYY',
          dateA11yLabel: 'LL',
          monthYearA11yLabel: 'MMMM YYYY'
        }
      });
    }
  }
}
