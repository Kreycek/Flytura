import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ConfigService } from '../../../services/config.service';



@Injectable({
  providedIn: 'root',
})
export class OutPutInvoicesService {
    // URL do seu backend para login

  constructor(
    private http: HttpClient,
    private configService:ConfigService
  ) {}


  
  getAllOutPutInvoicesDataPagination(page:number, limit:number, key?:string | null, companyCode?:string | null, startDate?:string | null, endDate?:string | null): Observable<any> {
    return this.http.get(
        this.configService.apiUrl + "/SearchOutPutInvoices?page="+page + 
        "&limit="+limit + 
        "&key="+key+ 
        "&companyCode="+companyCode+ 
        "&startDate="+startDate+ 
        "&endDate="+endDate , {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }


    addOutPutInvoices(data:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/InsertOutPutInvoices", data, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }


 
}