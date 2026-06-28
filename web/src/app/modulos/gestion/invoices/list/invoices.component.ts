import { CommonModule } from '@angular/common';
import { Component, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AirLineService } from '../../airLine/airLIne.service';
import { ConfigService } from '../../../../services/config.service';
import moment from 'moment';
import { ModalOkComponent } from '../../../../modal/modal-ok/modal-ok.component';
import { MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';
import { InvoicesService } from '../invoices.service';
import { PaginatorComponent } from '../../../../paginator/paginator.component';
import { jwtDecode } from 'jwt-decode';
import { EnumPerfil } from '../../../Enum/perfil';
import { ModalConfirmationComponent } from '../../../../modal/modal-confirmation/modal-confirmation.component';
import { ModuloService } from '../../../modulo.service';
import { SpinnerComponent } from '../../../../components/spinner/spinner.component';
import { AlertMoreColumnsComponent } from "../../../../components/alert-more-columns/alert-more-columns.component";
import { ImportInvoicesComponent } from '../import-invoices/import-invoices.component';
import { ModalMsgsComponent } from "../../../../modal/modal-msgs/modal-msgs.component";


@Component({
  selector: 'app-invoices', 
  imports: [
    CommonModule,
    FormsModule,
    MatDatepickerModule,
    TranslateModule,
    MatNativeDateModule,
    ModalOkComponent,
    PaginatorComponent,
    ModalConfirmationComponent,
    SpinnerComponent,
    AlertMoreColumnsComponent,
    ImportInvoicesComponent
    
],
  templateUrl: './invoices.component.html',
  styleUrl: './invoices.component.css',
   providers: [{ provide: MAT_DATE_LOCALE, useValue: 'pt-BR' }]
})
export class InvoicesComponent {

  @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  
  @ViewChild(ModalConfirmationComponent) modalConfirm!: ModalConfirmationComponent; 
    @ViewChild(ImportInvoicesComponent) modalImportInvoicesComponent!: ImportInvoicesComponent;
    searchStatusDonwload:string='';
    searchBilledFlytura:string='';
    searchAirlineCode:string='';
    searchAirlineDtInicio:string='';
    searchAirlineDtFim :string='';
    searchKey :string='';
    searchStatus:string='';
    statusImportData:any[]=[];
    airLInes:any[]=[];
    msgNotFound=false;
    totalRegistros: number = 0;
    totalPages: number = 1;    
    currentPage: number = 1;
    limit: number = 0;  
    dados:any;
    // imgsDownload:string[]=[]
    objPesquisar:any= {};
    adm:boolean=true;
    isLoading:boolean=false;
    viewMsgMoreColumns=false;
    FormData: FormData;
    decoded:any
    isSave:boolean=false;
    duplicateFiles:string[]=[];

    constructor(
            private airLineService: AirLineService,
            public configService:ConfigService,
            public invoicesService:InvoicesService,
            private translate: TranslateService,
            public moduloService:ModuloService,
          ) {}    

     

    ngOnInit() {

       const token = localStorage.getItem('token'); // ou onde você armazenou o JWT
            
        if (token) {
         this.decoded = jwtDecode<any>(token);
          if(this.decoded && this.decoded.perfis && this.decoded.perfis.length>0) {
              this.adm=this.decoded.perfis.some((data:any)=> data===EnumPerfil.ADM)
          }
        }

        this.limit=this.configService.limitPaginator;
        this.airLineService.getAllAirLine().subscribe((response)=>{
           this.airLInes=this.configService.sortByKey(response,'name');   
        
        });
        this.search(this.currentPage);           
    }

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

    async search(pageNumber:number) {       
          
        if(this.searchAirlineDtInicio || this.searchAirlineDtFim ) {
        
            if(await this.invalidDate(this.searchAirlineDtInicio, this.translate.instant('Ecra.invalidDateStart'), true)) {
                return false;
            }
  
            if(await this.invalidDate(this.searchAirlineDtFim, this.translate.instant('Ecra.invalidDateEnd'), true)) {
                return false;
            }
                
            const startDate = moment(this.searchAirlineDtInicio, 'YYYY-MM-DD');
            const endDate = moment(this.searchAirlineDtFim, 'YYYY-MM-DD');
  
            if (startDate.isAfter(endDate)) {
                const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.StartDateGreaterThanEndDate'),true);             
                    if (resultado) {
                      return false;
                      // Insira aqui a lógica para continuar após a confirmação
                    } else {
                      
                    }
              
            }
  
            if (endDate.isBefore(startDate)) {
              const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.StartDateEarlierThanEndDate'),true);             
                    if (resultado) {
                      return false;
                      // Insira aqui a lógica para continuar após a confirmação
                    } else {
                      
                    }          
              }
      }                       

      this.objPesquisar= {
            billedFlytura:this.searchBilledFlytura,
            key:this.searchKey,
            companyCode:this.searchAirlineCode,
            startDate:this.searchAirlineDtInicio  ? moment(this.searchAirlineDtInicio).format('YYYY-MM-DDTHH:mm:ss[Z]') : '',                               
            endDate: this.searchAirlineDtFim ? moment(this.searchAirlineDtFim).set({ hour: 23, minute: 59, second: 59 }).format('YYYY-MM-DDTHH:mm:ss[Z]')  : '' ,
            doDonwload :this.searchStatusDonwload,
      };
     
      this.invoicesService.getAllS3ImagesDBDataPagination(
                    pageNumber,
                    this.limit,
                    this.objPesquisar.billedFlytura,                     
                    this.objPesquisar.doDonwload,
                    this.objPesquisar.key,
                    this.objPesquisar.companyCode,
                    this.objPesquisar.startDate,
                    (this.objPesquisar.endDate && this.objPesquisar.endDate!=undefined) ? this.objPesquisar.endDate : ''
              ).subscribe((response:any)=>{
                // console.log('response.imagesDB',response.imagesDB);
                if(response.imagesDB) {
                    this.dados=response.imagesDB.map((element:any) => {
                        element.PasteName=element.FileName?.replace(".zip", "");
                        return element
                    });
                } else {
                   this.dados=null;
                }

                  this.totalRegistros = response.total;

                    if(this.totalRegistros==0) {
                     this.currentPage=1;
                    }

                  this.totalPages = response.pages;
              })

        return null;
    }

     
    onPageChange(newPage: number) {
      this.currentPage = newPage;
      this.search(this.currentPage);
    }

     donwloadAll() {
      // console.log('this.searchAirlineDtInicio ',this.searchAirlineDtInicio);
      // console.log('this.searchAirlineDtFim ',this.searchAirlineDtFim);
      this.objPesquisar={
              companyCode:this.searchAirlineCode,
              startDate:this.searchAirlineDtInicio  ? moment(this.searchAirlineDtInicio).format('YYYY-MM-DDTHH:mm:ss[Z]') : '',                               
              endDate: this.searchAirlineDtFim ? moment(this.searchAirlineDtFim).set({ hour: 23, minute: 59, second: 59 }).format('YYYY-MM-DDTHH:mm:ss[Z]')  : ''     
      }

      this.invoicesService.getAllS3ImagesDBFull(this.objPesquisar.companyCode,this.objPesquisar.startDate,this.objPesquisar.endDate)
         .subscribe(async (response:any)=>{


           if (response.imagesDB && response.imagesDB.length > 0) {
            //  console.log('response.imagesDB',response.imagesDB);

                  this.isLoading=true;                  
                  let ids:String[]=[]

                  interface ImageResponse {
                    ID:string
                    ZipFileName: string;
                    DownloadDone:boolean
                    // outras propriedades, se houver
                  }
                  
                  for (const item of response.imagesDB as ImageResponse[]) {
                    // this.imgsDownload.push(item.ZipFileName);
                    ids.push(item.ID)
                  }
                     
                  try {
                      const ok = await this.invoicesService.downloadGroupedZipAndReturnTrue(response.imagesDB);
                      if (ok) {
                          this.isLoading=false;
                        // chegou ao fim com sucesso
                        // console.log('Concluído e retornou true');
                        // opcional: mostrar toast/snackbar
                      }
                    } catch (e) {
                      this.isLoading=false;
                        const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.downloadInvoicesError'), true);
                        if (resultado) {                               
                          
                          return; // apenas interrompe a execução
                        }
                    }

                    this.invoicesService.updateMultipleStatusS3Images({
                      ids:ids,
                      DownloadDone:true
                    }).subscribe((response:any)=>{
                      for (const item of this.dados as ImageResponse[]) {
                        item.DownloadDone=true;
                      }
                    })
                  }        
            })
     }

    donwloadJustOne(linha:any) {
      const objUpdate={Id:linha.ID, DownloadDone:true}
      linha.DownloadDone=true;
      this.invoicesService.updateStatusS3Image(objUpdate).subscribe()
    }


  updateStatusDownloadPdfOrXml(linha:any,fileType:string) {
      if(fileType=='pdf')
        linha.DownloadPDFDone=true;

      if(fileType=='xml')
         linha.DownloadXMLDone=true;

      const objUpdate={id:linha.ID, fileType:fileType,downloadOk:true}
      this.invoicesService.UpdateStatusPdforXml(objUpdate).subscribe()
    }

    
    async deleteItem(id:string, filenName:string, zipFileName:string, pdfFileName:string, xmlFileName:string ) {
        if(!zipFileName) {
          zipFileName=filenName
        }
        const resultado = await this.modalConfirm.openModal(
                  true,
                  this.translate.instant('Ecra.deleteInformation'),
                  this.translate.instant('Ecra.yes'),
                  this.translate.instant('Ecra.no')
                ); 
        if (resultado) {    
          
              this.modalConfirm.isVisible=false;
              this.invoicesService.deleteS3Images(id,zipFileName,pdfFileName,xmlFileName).subscribe((trt)=>{
                this.search(this.currentPage);      
              });
               
        } else {
          this.modalConfirm.isVisible=false;      
        }
    }     

    async importManual() {
      
       const resultado = await this.modalImportInvoicesComponent.openModal(
                  true,
                  this.translate.instant('Ecra.deleteInformation')
                  ,this.airLInes
                ); 
        if (resultado) {     
          
        }
    }


    recieveFormData(formData: FormData) {   

      formData.append('idUserInserted', this.decoded.idUser);
      formData.append('userName', this.decoded.name);

      

     const filesArray: string[] = [];

    formData.forEach((value, key) => {
      if (value instanceof File) {
        filesArray.push(value.name); // ✅ pega o nome do arquivo
      }
    });

    console.log(filesArray);

      this.invoicesService.checkDuplicatePDFsHandler(filesArray,"/CheckDuplicatePDFs").subscribe((response:any)=>{
        //Para poder ativar o ngOnChange dentro do componente de modal faz isso abaixo
        this.isSave=false;

        if(response && response.duplicates && response.duplicates.length>0) {
          this.duplicateFiles=response.duplicates;
        } 
        else {
          this.moduloService.importFile(formData, "/UploadManualImportInvoicesRecord").subscribe({
            next: async (returnSheet: any) => {
              if (returnSheet.message) {
                 this.isSave=true;
               
                 this.search(1); 
              }
            },
            error: async (err: any) => {
              console.error('Erro do servidor:', err);
              if (err.error?.status === "500") {
                
              }
              else if (err.status === 500) {
                 
                }
              }
            
          });
        }
      })

        
     
    }

        
    setFiles(filesString: string) {
      return filesString ? filesString.split(';') : [];
    }


    ensurePdfExtension(fileName: string): string {
      if (!fileName) return '';

      return fileName.toLowerCase().endsWith('.pdf')
        ? fileName
        : fileName + '.pdf';
    }

    getFileNameFromUrl(url: string): string {
      if (!url) return '';

      return url.split('/').pop() || '';
    }
     
}
