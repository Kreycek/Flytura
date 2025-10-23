import { CommonModule } from '@angular/common';
import { Component, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { TranslateModule } from '@ngx-translate/core';
import { AirLineService } from '../../companys/airLine/airLIne.service';
import { ConfigService } from '../../../services/config.service';
import moment from 'moment';
import { ModalOkComponent } from '../../../modal/modal-ok/modal-ok.component';
import { MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';
import { InvoicesService } from './invoices.service';
import { PaginatorComponent } from '../../../paginator/paginator.component';

@Component({
  selector: 'app-invoices', 
  imports: [CommonModule,FormsModule,MatDatepickerModule,TranslateModule,MatNativeDateModule,ModalOkComponent,PaginatorComponent],
  templateUrl: './invoices.component.html',
  styleUrl: './invoices.component.css',
   providers: [{ provide: MAT_DATE_LOCALE, useValue: 'pt-BR' }]
})
export class InvoicesComponent {

  @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  
   searchAirlineCode:string='';
    searchAirlineDtInicio:string='';
    searchAirlineDtFim :string='';
    searchStatus:string='';
    statusImportData:any[]=[];
    airLInes:any[]=[];
    msgNotFound=false;
    totalRegistros: number = 0;
    totalPages: number = 1;    
    currentPage: number = 1;
    limit: number = 0;  
    dados:any;
    imgsDownload:string[]=[]
     objPesquisar:any= { 
          
         
      };
    
    /**
     *
     */
    constructor(
      private airLineService: AirLineService,
            public configService:ConfigService,
            public invoicesService:InvoicesService
          ) {}

     async invalidDate(date: string, msg:string, showAlert:boolean=false) : Promise<boolean> {              
              const value = moment(date); // inputDate pode ser string, Date, etc.
      
              if (!value.isValid()) {
                if(showAlert) {
                    const resultado = await this.modalOk.openModal(msg,true);             
                        if (resultado) {
                      
                          // Insira aqui a lógica para continuar após a confirmação
                        } else {
                          
                        }
                    }
                      
                return true;
              } else {
               return false;
              }
    }

    ngOnInit() {

        this.limit=this.configService.limitPaginator;
        this.airLineService.getAllAirLine().subscribe((response)=>{
          this.airLInes=response;
        });
        this.search(this.currentPage);
           
    }

    async search(pageNumber:number) {

       
          
        if(this.searchAirlineDtInicio || this.searchAirlineDtFim ) {
        
            if(await this.invalidDate(this.searchAirlineDtInicio, 'Data de início inválida, ou se a data de fim estiver preenchida é necessário preencher uma data de início.', true)) {
                return false;
            }
  
            if(await this.invalidDate(this.searchAirlineDtFim, 'Data de fim inválida, ou se a data de início estiver preenchida é necessário preencher uma data de fim.', true)) {
                return false;
            }
                
            const startDate = moment(this.searchAirlineDtInicio, 'YYYY-MM-DD');
            const endDate = moment(this.searchAirlineDtFim, 'YYYY-MM-DD');
  
            if (startDate.isAfter(endDate)) {
                const resultado = await this.modalOk.openModal('Data de início não pode ser maior que a data de fim',true);             
                    if (resultado) {
                      return false;
                      // Insira aqui a lógica para continuar após a confirmação
                    } else {
                      
                    }
              
            }
  
            if (endDate.isBefore(startDate)) {
              const resultado = await this.modalOk.openModal('Data de fim é menor que a data de inicio',true);             
                    if (resultado) {
                      return false;
                      // Insira aqui a lógica para continuar após a confirmação
                    } else {
                      
                    }          
              }
      }                       

      this.objPesquisar=   {companyCode:this.searchAirlineCode,
            startDate:this.searchAirlineDtInicio  ? moment(this.searchAirlineDtInicio).format('YYYY-MM-DDTHH:mm:ss[Z]') : '',                               
            endDate: this.searchAirlineDtFim ? moment(this.searchAirlineDtFim).set({ hour: 23, minute: 59, second: 59 }).format('YYYY-MM-DDTHH:mm:ss[Z]')  : ''     }
      this.invoicesService.getAllS3ImagesDBDataPagination(
        pageNumber,
        this.limit,
        this.objPesquisar.companyCode,
        this.objPesquisar.startDate,
        this.objPesquisar.endDate).subscribe((response:any)=>{
       
  console.log('Faturas',response);
      
               this.dados=response.imagesDB;    
          this.totalRegistros = response.total;
          this.totalPages = response.pages;
      })

      return null;
    }

     
    onPageChange(newPage: number) {
      this.currentPage = newPage;
      this.search(this.currentPage);
    }

     donwloadAll() {
      this.objPesquisar={
              companyCode:this.searchAirlineCode,
              startDate:this.searchAirlineDtInicio  ? moment(this.searchAirlineDtInicio).format('YYYY-MM-DDTHH:mm:ss[Z]') : null,                               
              endDate: this.searchAirlineDtFim ? moment(this.searchAirlineDtFim).set({ hour: 23, minute: 59, second: 59 }).format('YYYY-MM-DDTHH:mm:ss[Z]')  : null     
      }

      this.invoicesService.getAllS3ImagesDBFull(this.objPesquisar.companyCode,this.objPesquisar.startDate,this.objPesquisar.endDate)
         .subscribe((response:any)=>{

                  this.imgsDownload = [];
                  let ids:String[]=[]

                  interface ImageResponse {
                    ID:string
                    ZipFileName: string;
                    DownloadDone:boolean
                    // outras propriedades, se houver
                  }

                  for (const item of response.imagesDB as ImageResponse[]) {
                    this.imgsDownload.push(item.ZipFileName);
                    ids.push(item.ID)
                  }

          
                  if (this.imgsDownload.length > 0) {
                    console.log('this.imgsDownload ',this.imgsDownload);
                      this.invoicesService.downloadZip(this.imgsDownload);
                  }

                  this.dados=response.imagesDB;    

                   this.invoicesService.updateMultipleStatusS3Images({
                    ids:ids,
                    DownloadDone:true
                   }).subscribe((response:any)=>{
                     for (const item of this.dados as ImageResponse[]) {
                      item.DownloadDone=true;
                  }
                   })
        
      })


     }

    donwloadJustOne(linha:any) {

      const objUpdate={Id:linha.ID, DownloadDone:true}

      linha.DownloadDone=true;

      this.invoicesService.updateStatusS3Image(objUpdate).subscribe()

    }
}
