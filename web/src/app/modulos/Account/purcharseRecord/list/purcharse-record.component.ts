
import { Component, ElementRef, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ConfigService } from '../../../../services/config.service';
import { ModalConfirmationComponent } from '../../../../modal/modal-confirmation/modal-confirmation.component';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { PaginatorComponent } from '../../../../paginator/paginator.component';
import * as _moment from 'moment';
import { PurcharseRecordService } from '../purcharse-record.service';
import { ModalOkComponent } from '../../../../modal/modal-ok/modal-ok.component';
import { SpinnerComponent } from '../../../../components/spinner/spinner.component'
import { AirLineService } from '../../../gestion/airLine/airLIne.service';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';
import moment from 'moment';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ModelsComponent } from '../../../gestion/sheetModels/models.component';
import { ModuloService } from '../../../modulo.service';
import { ModalMsgsComponent } from '../../../../modal/modal-msgs/modal-msgs.component';
import { jwtDecode } from 'jwt-decode';
import { EnumPerfil } from '../../../Enum/perfil';
import { AlertMoreColumnsComponent } from "../../../../components/alert-more-columns/alert-more-columns.component";

@Component({
  selector: 'app-invoices',
  imports: [
    CommonModule,
    FormsModule,
    PaginatorComponent,
    ModalOkComponent,
    MatDatepickerModule,
    MatNativeDateModule,
    TranslateModule,
    ModelsComponent,
    ModalMsgsComponent,
    ModalConfirmationComponent,
    SpinnerComponent,
    AlertMoreColumnsComponent
],
  templateUrl: './purcharse-record.component.html',
  styleUrl: './purcharse-record.component.css',
  providers: [{ provide: MAT_DATE_LOCALE, useValue: 'pt-BR' }],
})
export class PurcharseRecordComponent { 
      @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  
      @ViewChild(ModelsComponent) modalDocuments!: ModelsComponent;
      @ViewChild(ModalMsgsComponent) modalMsgsComponent!: ModalMsgsComponent;
      @ViewChild('fileInput') fileInput!: ElementRef;
      @ViewChild(ModalConfirmationComponent) modalConfirm!: ModalConfirmationComponent;

      searchKey: any = '';
      searchName: any =  '';
      searchLastName: any = '';
      searchAirlineCode:any = '';
      searchAirlineDtInicio:any = '';
      searchAirlineDtFim :any = '';
      searchStatus:any = '';
      statusMessage: any =  '';
      totalRegistros: number = 0;
      totalPages: number = 1;
      currentPage: number = 1;
      limit: number = 0;  
      currentYear: number = new Date().getFullYear();
      dados:any
      statusImportData:any[]=[]
      costCenterSubModalList:any[]=[]
      costCenters:any[]=[]    
      costCentersSub:any[]=[]
      airLInes:any[]=[]
      adm:boolean=true;
      isLoading:boolean=false;
      fillOneFilter=false;   

    constructor(
      private router: Router, 
      private purcharseRecordService: PurcharseRecordService,
      private airLineService: AirLineService,
      public configService:ConfigService,
      public moduloService:ModuloService,
      private translate: TranslateService,
      private route: ActivatedRoute,
      
    ) {} 

   decoded:any

    ngOnInit() {    

        const token = localStorage.getItem('token'); // ou onde você armazenou o JWT        
        if (token) {
          this.decoded = jwtDecode<any>(token);
          if( this.decoded &&  this.decoded.perfis &&  this.decoded.perfis.length>0) {
              this.adm= this.decoded.perfis.some((data:any)=> data===EnumPerfil.ADM)
          }
        }

        this.limit=this.configService.limitPaginator;

        this.airLineService.getAllAirLine().subscribe((response)=>{    
          this.airLInes=this.configService.sortByKey(response,'name');  

        });

        this.moduloService.getAllStatusImportData().subscribe((response:any)=>{     
          this.statusImportData=response;
        });

        //Essa parte foi implementada para caso edite passa os dados para 
        // edição e se retornar pega novamene os parâmetros mantendo a busca
        this.route.paramMap.subscribe(params => {
          this.searchKey=params?.get('key') ?? ''
          this.searchName=params?.get('nombre') ?? ''
          this.searchLastName=params.get('apellido') ?? ''
          this.searchAirlineCode=params.get('companyCodAirline') ?? ''
          this.searchAirlineDtInicio=params.get('dtStart') ?? ''
          this.searchAirlineDtFim=params.get('dtEnd') ?? ''
          this.searchStatus=params.get('status') ?? ''

          if(
          !this.searchKey && 
          !this.searchName && 
          !this.searchLastName && 
          !this.searchAirlineCode && 
          !this.searchAirlineDtInicio&& 
          !this.searchAirlineDtFim&& 
          !this.searchStatus) {
            this.loadDataAfterImport();
          }
          else {
              this.searchPurcharseRecord(this.currentPage);
          }
        }) 
  }


  loadDataAfterImport() {
        this.purcharseRecordService.getAllPurchaseRecordDataPagination(this.currentPage,this.limit).subscribe(async (response:any)=>{                 
            this.dados=response.purcharseRecord;  
            this.totalRegistros = response.total;
            this.totalPages = response.pages;              
        });
       
  }
      
  onFileSelected(event: any) {
      const file: File = event.target.files[0];
     
      if (file) {
              const formData = new FormData();
              formData.append('file', file);
              formData.append('idUserInserted', this.decoded.idUser);
              this.isLoading=true;

          this.moduloService.importFile(formData, "/UploadExcelPurcharseRecord").subscribe({
            next: async (returnSheet: any) => {
              if (returnSheet.message) {
                // Caso 1: Planilha vazia
                if (returnSheet.message.emptySheet) {
                    this.isLoading=false;
                  const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.sheetEmpty'), true);
                  if (resultado) {
                    this.fileInput.nativeElement.value = '';
                    
                    return; // apenas interrompe a execução
                  }
                }
                else  if (returnSheet.message.ivalidNumberCols) {
                  this.isLoading=false;
                  const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.ivalidNumberCols'), true);
                  if (resultado) { 
                     this.fileInput.nativeElement.value = '';
                    return; // apenas interrompe a execução
                  }
                }
                // Caso 2: Erros específicos
                else if (
                  returnSheet.message.emptyError ||
                  returnSheet.message.cientificError ||
                  returnSheet.message.ivalidFormatError
                ) {
                  const errorAlert: string[] = [];

                  if (returnSheet.message.emptyError) {
                     this.isLoading=false;
                    errorAlert.push(
                      this.translate.instant('Ecra.emptySheetLinesStart') +
                        returnSheet.message.emptyError +
                        this.translate.instant('Ecra.emptySheetLinesEnd') +
                        '<br/><br/>'
                    );
                  }

                  if (returnSheet.message.cientificError) {
                     this.isLoading=false;
                    errorAlert.push(
                      this.translate.instant('Ecra.cientificSheetLinesStart') +
                        returnSheet.message.cientificError +
                        this.translate.instant('Ecra.cientificSheetLinesEnd') +
                        '<br/><br/>'
                    );
                  }

                  if (returnSheet.message.ivalidFormatError) {
                     this.isLoading=false;
                    errorAlert.push(
                      this.translate.instant('Ecra.ivalidFormatErrorSheetLineStart') +
                        returnSheet.message.ivalidFormatError +
                        this.translate.instant('Ecra.ivalidFormatErrorSheetLineEnd') +
                        '<br/><br/>'
                    );
                  }

                  const resultado = await this.modalMsgsComponent.openModal(errorAlert, '', true);
                  if (resultado) {
                     this.isLoading=false;
                    this.fileInput.nativeElement.value = '';
                    
                    return;
                  }
                }
                // Caso 3: Nenhum registro importado
                else if (returnSheet.message.totalEmpty === 0 && returnSheet.message.totalRecordsImport === 0) {
                    this.isLoading=false;
                  const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.sheetImportedAfter'), true);
                  if (resultado) {
                 
                    this.loadDataAfterImport();
                    this.fileInput.nativeElement.value = '';
                    return;
                  }
                }
                
                // Caso 4: Importação parcial
                else if (returnSheet.message.totalEmpty > 0) {
                   this.isLoading=false;
                  const resultado = await this.modalOk.openModal(
                    this.translate.instant('Ecra.sheetImportedProblemPart1') +
                      returnSheet.message.totalEmpty +
                      this.translate.instant('Ecra.sheetImportedProblemPart2'),
                    true
                  );
                  if (resultado) {                      
                    this.loadDataAfterImport();
                    this.fileInput.nativeElement.value = '';
                    return;
                  }
                }
                // Caso 5: Importação OK
                else {
                   this.isLoading=false;
                  this.loadDataAfterImport();
                  const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.sheetImportedOk'), true);
                  if (resultado) {
                    this.fileInput.nativeElement.value = '';
                  }
                }
              }
            },
            error: async (err: any) => {
              console.error('Erro do servidor:', err);
              if (err.error?.status === "500") {
                 this.isLoading=false;
                const resultado = await this.modalOk.openModal(                
                   this.translate.instant('Ecra.correctSheetName')
                   .replace("{0}", err.error?.codFile)
                   .replace("{1}", err.error?.sheetName)
                   .replace("{2}", err.error?.codFile)
                   .replace("{3}", err.error?.nameAirLine)
                   .replace("{4}", err.error?.codFile)
                    .replace("{5}", err.error?.fileName)
                  , true);
                if (resultado) {
                    
                  this.fileInput.nativeElement.value = '';
                }
              }
              else if (err.status === 500) {
                   this.isLoading=false;
                const resultado = await this.modalOk.openModal(                
                    this.translate.instant('Ecra.'+err.error.replace("\n","")), 
                    true);
                     if (resultado) {
                      
                  this.fileInput.nativeElement.value = '';
                }
              }
            }
          });
        } 
    }
    
  
    onPageChange(newPage: number) {
      this.currentPage = newPage;
      this.searchPurcharseRecord(this.currentPage);
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
    
    async searchPurcharseRecord(currentPage:number) {    
  
      let objPesquisar: { 
          key: string;
          name: string;   
          lastName: string;    
          companyCode:string;  
          startDate?:string | null;
          endDate?:string | null;
          status?:string | null;
          statusText?:string | null;
          page:number;
          limit:number;
      }

      if(this.searchAirlineDtInicio || this.searchAirlineDtFim ){
     
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

  
      objPesquisar= { 
        key: this.searchKey, 
        name: this.searchName, 
        lastName: this.searchLastName, 
        companyCode:this.searchAirlineCode,
        startDate:this.searchAirlineDtInicio ? moment(new Date(this.searchAirlineDtInicio)).format('YYYY-MM-DDT00:00:00Z') : null,
        endDate:this.searchAirlineDtFim ? moment(new Date(this.searchAirlineDtFim)).format('YYYY-MM-DDT23:59:59Z') :null,
        status:this.searchStatus,
        statusText:this.statusMessage,        
        page:currentPage,
        limit:this.limit
      };
      
     
      this.purcharseRecordService.searchPurchaseRecordData(objPesquisar).subscribe((response:any)=>{
        this.dados=response.purcharseRecord;    
        this.totalRegistros = response.total;
        if(this.totalRegistros==0) {
          this.currentPage =1;
        }
        this.totalPages = response.pages;
      })
      return null;  
    }
  
    
    async searchPurcharseRecordByPeriod(currentPage:number) {    
  
      let objPesquisar: { 
          key: string;
          name: string;   
          lastName: string;    
          companyCode:string;  
          startDate?:string | null;
          endDate?:string | null;
          status?:string | null;
          statusText?:string | null;
      }

      if(this.searchAirlineDtInicio || this.searchAirlineDtFim ){
     
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

  
      objPesquisar= { 
        key: this.searchKey, 
        name: this.searchName, 
        lastName: this.searchLastName, 
        companyCode:this.searchAirlineCode,
        startDate:this.searchAirlineDtInicio ? moment(new Date(this.searchAirlineDtInicio)).format('YYYY-MM-DDT00:00:00Z') : null,
        endDate:this.searchAirlineDtFim ? moment(new Date(this.searchAirlineDtFim)).format('YYYY-MM-DDT23:59:59Z') :null,
        status:this.searchStatus,
        statusText:this.statusMessage
      };
      
     
      this.purcharseRecordService.searchPurchaseRecordDataByPeriod(objPesquisar).subscribe((response:any)=>{
        this.dados=response.purcharseRecord;            
      })

      return null;
  
  }

  addPurcharseRecord() {
      this.router.navigate(['/aplicacao/addPurchaseRecord']);
  }
  
  updatePurcharseRecord(id:string) {
      const dtStart =this.searchAirlineDtInicio ? new Date(this.searchAirlineDtInicio).toISOString() : '';
      const dtEnd =this.searchAirlineDtInicio ? new Date(this.searchAirlineDtFim).toISOString() : '';
      this.router.navigate(['/aplicacao/addPurchaseRecord', id, this.searchKey,this.searchName,this.searchLastName,this.searchAirlineCode,dtStart,dtEnd,this.searchStatus]);   
  }   
  
  tooltip = {
    visible: false,
    text: '',
    top: 0,
    left: 0,
    maxWidth: 'auto'
  };

  mousemoveHandler(event: MouseEvent, item: any, show:boolean=true, textMessage:string) {
    const offsetX = 15;
    const offsetY = 15;
    const tooltipWidth = 200;
    const screenWidth = window.innerWidth;

    let left = event.clientX + offsetX;
    let top = event.clientY + offsetY;
    let maxWidth = 'auto';

    if (left + tooltipWidth > screenWidth) {
      left = screenWidth - tooltipWidth - offsetX;
      maxWidth = '180px';
    }

    this.tooltip = {
      visible: show,
      text:textMessage,
      top,
      left,
      maxWidth
    };
  }

    hideCustomTooltip() {
      this.tooltip.visible = false;
    }


    async openModels() {
      const resultado = await this.modalDocuments.openModal(
        [],
        "Lista de documentos da empresa <br\><br\>",
          true); 

        if (resultado) {      
        } else {
          
        }
    }  

    
    async deleteItem(id:string) {

        const resultado = await this.modalConfirm.openModal(
                  true,
                  this.translate.instant('Ecra.deleteInformation'),
                  this.translate.instant('Ecra.yes'),
                  this.translate.instant('Ecra.no')
                ); 
        if (resultado) {     
          
              this.modalConfirm.isVisible=false;
              this.purcharseRecordService.deletePurcharseRecord(id).subscribe((response)=>{
                 this.searchPurcharseRecord(this.currentPage); 
              });
             
        } else {
          this.modalConfirm.isVisible=false;
      
        }
    }   
    
    clearSearchFields() {
      this.searchKey=''
      this.searchName=''
      this.searchLastName=''
      this.searchAirlineCode=''
      this.searchAirlineDtInicio= ''
      this.searchAirlineDtFim= ''
      this.searchStatus= ''
    }
  
    async generateReportExcel() {


        if(this.searchAirlineDtInicio || this.searchAirlineDtFim ){
     
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

       const objPesquisar= { 
        key: this.searchKey, 
        name: this.searchName, 
        lastName: this.searchLastName, 
        companyCode:this.searchAirlineCode,
        startDate:this.searchAirlineDtInicio ? moment(new Date(this.searchAirlineDtInicio)).format('YYYY-MM-DDT00:00:00Z') : null,
        endDate:this.searchAirlineDtFim ? moment(new Date(this.searchAirlineDtFim)).format('YYYY-MM-DDT23:59:59Z') :null,
        status:this.searchStatus,
        statusText:this.statusMessage
      };
        this.purcharseRecordService.searchPurchaseRecordDataByPeriod(objPesquisar).subscribe((response:any)=>{
          this.dados=response.purcharseRecord;              
           this.moduloService.exportToExcelPurcharseRecord(this.dados,'teste')
        })

        return null
    }

    verifyFieldsSearchFill() {   
      
      const campos = [
        this.searchKey,
        this.searchName,
        this.searchLastName,
        this.searchAirlineCode,
        this.searchAirlineDtInicio,
        this.searchAirlineDtFim,
        this.searchStatus
      ];

      this.fillOneFilter = campos.some(v => v !== null && v != '');    

    }
}
