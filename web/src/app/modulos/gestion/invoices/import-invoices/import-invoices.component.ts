import { CommonModule } from '@angular/common';
import {Component, ElementRef, input, output, SimpleChanges, ViewChild } from '@angular/core';
import {  FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { Subject } from 'rxjs';
import { AirLineService } from '../../airLine/airLIne.service';
import { ConfigService } from '../../../../services/config.service';
import { ModalOkComponent } from '../../../../modal/modal-ok/modal-ok.component';

@Component({
  selector: 'app-import-invoices',
  imports: [CommonModule, ReactiveFormsModule, FormsModule, TranslateModule, TranslateModule, ModalOkComponent],
  templateUrl: './import-invoices.component.html',
  styleUrl: './import-invoices.component.css'
})
export class ImportInvoicesComponent {

  isVisible = false;
  message: string = '';
  documentMiniFormDescription:string=''
  airLInes:any[]=[]  
  formulario: FormGroup| null = null;
  filesForImport: File[] = [];
  formData = new FormData();
  formDataOutput = output<FormData>();
  msgGravar = input<boolean>(false);
  @ViewChild('fileInput') fileInput!: ElementRef;
  private responseSubject = new Subject<boolean>();  
  @ViewChild(ModalOkComponent) modal!: ModalOkComponent;   
  

  constructor( private fb: FormBuilder, 
    private translate: TranslateService,
     private airLineService: AirLineService,
     public configService:ConfigService,
  ) {}

  ngOnInit() {
    this.createForm({CompanyCode:"",Active:true})     
  }  

  ngOnChanges(changes: SimpleChanges) {  
      if (changes['msgGravar']) {
          if(changes['msgGravar'].currentValue) {
            this.handleMsgGravar()
              this.fileInput.nativeElement.value = ''
              this.formData = new FormData();
              this.filesForImport=[];
               this.createForm({CompanyCode:"",Active:true}) ;
          }
      }
  }
  // Método para abrir o modal e retornar um Observable
  openModal(
    isVisible:boolean,
     message:string, 
     airLInes:any[]=[],    
    
    ): Promise<boolean> {
    this.isVisible = isVisible;
    this.message=message;   
    this.airLInes=airLInes
    return new Promise(resolve => {
      this.responseSubject = new Subject<boolean>();
      this.responseSubject.subscribe(response => {
        // this.isVisible = false;
        resolve(response);
      });
    });
  }

  confirm() {
    
    this.responseSubject.next(true);
    this.responseSubject.complete();
  }

  cancel() {
    this.responseSubject.next(false);
    this.responseSubject.complete();
    this.isVisible =false;
  }

   createForm(obj:any) {
          this.formulario = this.fb.group({
            active: [obj.Active],
            companyCode: [obj.CompanyCode, Validators.required],
            key: [obj.Key, Validators.required],                  
            flyturaInvoice: [obj.FlyturaInvoice??''],          
          });              
      }      
          

    fillFormData(files:File[], companyCode:string, key:string, billedFlytura:string) {
      files.forEach(file => {
        this.formData.append('file', file);
      });   

      this.formData.append('companyCode', companyCode);
      this.formData.append('key', key);
      this.formData.append('billedFlytura', billedFlytura); // S ou N
    }

    onFileSelected(event: Event) {

          const input = event.target as HTMLInputElement;
          if (input.files && input.files.length > 0) {          
            this.filesForImport=Array.from(input.files);
             this.fillFormData(
              this.filesForImport,
              this.formulario?.controls["companyCode"].value,
              this.formulario?.controls["key"].value,
              this.formulario?.controls["flyturaInvoice"].value
            )          
          }
    }   
    
    removeFile(index: number) {
        
      this.filesForImport.splice(index, 1);
    
      this.fillFormData(
              this.filesForImport,
              this.formulario?.controls["companyCode"].value,
              this.formulario?.controls["key"].value,
              this.formulario?.controls["flyturaInvoice"].value
            )

    }

    async gravar() {    

       if (this.formulario?.invalid) {
         
        this.formulario.markAllAsTouched();
        return;
      }
      
      this.formDataOutput.emit(this.formData)     
         
    }   

    async handleMsgGravar() {
      if (this.msgGravar()) {
        const resultado = await this.modal.openModal(
          this.translate.instant('Ecra.purcharseRecordAddSuccess'),
          true
        );

        if (resultado) {
          // lógica após confirmação
        } else {
          // lógica cancelamento (se houver)
        }
      }
  }
}
