
  import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ConfigService } from '../../../services/config.service';



@Injectable({
  providedIn: 'root',
})
export class  ConciliationService {
    // URL do seu backend para login

      constructor(
    private http: HttpClient,
    private configService:ConfigService
  ) {}

  /*
Função criada por Ricardo Silva Ferreira
Inicio da criação 22/05/2026 15:38
Data Final da criação :  22/05/2026 15:38
*/

   getAllConciliationDataPagination(
    page:number, 
    limit:number, 
    originLocator?:string | null, 
    returnLocator?:string | null, 
    originETicket?:string | null, 
    returnETicket?:string | null, 
    startDate?:string | null, 
    endDate?:string | null): Observable<any> {
    return this.http.get(
        this.configService.apiUrl + "/SearchConciliationPagination?page="+page + 
        "&limit="+limit + 
        "&originLocator="+originLocator+ 
        "&returnLocator="+returnLocator+ 
        "&originETicket="+originETicket+ 
        "&returnETicket="+returnETicket+ 
        "&startDate="+startDate+ 
        "&endDate="+endDate , {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }


  
  /*
Função criada por Ricardo Silva Ferreira
Inicio da criação 22/05/2026 15:35
Data Final da criação :  22/05/2026 15:40
*/
  

getAllConciliationExcel(
  originLocator?: string | null, 
  returnLocator?: string, 
  originETicket?: string, 
  returnETicket?: string, 
  startDate?: string | null, 
  endDate?: string | null
): Observable<any> {

  let params = new HttpParams();

  console.log(originLocator,returnLocator,originETicket,returnETicket,)

  if (originLocator) params = params.set('originLocator', originLocator);
  if (returnLocator) params = params.set('returnLocator', returnLocator);
  if (originETicket) params = params.set('originETicket', originETicket);
  if (returnETicket) params = params.set('returnETicket', returnETicket);
  if (startDate) params = params.set('startDate', startDate);
  if (endDate) params = params.set('endDate', endDate);

  return this.http.get(
    this.configService.apiUrl + "/SearchConciliationExcel",
    {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
      params: params
    }
  );
}
}
