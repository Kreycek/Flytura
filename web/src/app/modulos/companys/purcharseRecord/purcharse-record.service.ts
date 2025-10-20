import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ConfigService } from '../../../services/config.service';

@Injectable({
  providedIn: 'root',
})
export class PurcharseRecordService {
    // URL do seu backend para login

  constructor(
    private http: HttpClient,
    private configService:ConfigService
  ) {}

  getAllPurchaseRecordDataPagination(page:number, limit:number): Observable<any> {
    return this.http.get(this.configService.apiUrl + "/GetAllPurcharseRecordPagination?page="+page + "&limit="+limit , {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

    //Carregar todos os diários sem documentos
    getAllPurchaseRecordData(): Observable<any> {
      return this.http.get(this.configService.apiUrl + "/GetAllPurcharseRecord" , {
        headers: new HttpHeaders({
          'Content-Type': 'application/json',
        }),
      });
    }
  

  getPurchaseRecordDataById(id:string): Observable<any> {
    return this.http.get(this.configService.apiUrl + "/GetPurcharseRecordById?id="+id, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  addPurchaseRecordData(data:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/InsertPurcharseRecord", data, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  updatePurchaseRecordData(data:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/UpdatePurcharseRecord", data, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  verifyExistPurchaseRecordData(data:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/VerifyExistPurcharseRecord", data, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  searchPurchaseRecordData(data:any): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/SearchPurcharseRecord", data, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    });
  }

  
  importarPlanilha(formData:FormData): Observable<any> {
    return this.http.post(this.configService.apiUrl + "/UploadPurcharseRecord", formData);

  }


    
  //Carregar todos os status de importação
   
GroupByCompanyName(
  status?: string | null,
  companyName?: string| null,
  startDate?: string| null,
  endDate?: string| null
): Observable<any> {
  let params = new HttpParams();

  if (status) {
    params = params.set('status', status);
  }
  if (companyName) {
    params = params.set('companyName', companyName);
  }
  if (startDate) {
    params = params.set('startDate', startDate);
  }
  if (endDate) {
    params = params.set('endDate', endDate);
  }

  return this.http.get(this.configService.apiUrl + '/GroupByCompanyName', {
    headers: new HttpHeaders({
      'Content-Type': 'application/json',
    }),
    params: params
  });
}

}
