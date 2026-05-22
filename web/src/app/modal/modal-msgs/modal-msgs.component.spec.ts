import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModalMsgsComponent } from './modal-msgs.component';

describe('ModalMsgsComponent', () => {
  let component: ModalMsgsComponent;
  let fixture: ComponentFixture<ModalMsgsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ModalMsgsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ModalMsgsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
