import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModalOkComponent } from './modal-ok.component';

describe('ModalOkComponent', () => {
  let component: ModalOkComponent;
  let fixture: ComponentFixture<ModalOkComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ModalOkComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ModalOkComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
