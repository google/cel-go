import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';

import { ReplApiService } from './repl-api.service';

describe('ReplApiService', () => {
  let service: ReplApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule]
    });
    service = TestBed.inject(ReplApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
