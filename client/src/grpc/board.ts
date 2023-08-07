/* eslint-disable */
import {
  CallOptions,
  ChannelCredentials,
  Client,
  ClientOptions,
  ClientUnaryCall,
  handleUnaryCall,
  makeGenericClientConstructor,
  Metadata,
  ServiceError,
  UntypedServiceImplementation,
} from "@grpc/grpc-js";
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "board";

export interface Likes {
  userId: string;
  questionId: number;
}

export interface NewSubject {
  title: string;
}

export interface Subject {
  id: number;
  title: string;
}

export interface SubjectId {
  id: number;
}

export interface NewQuestion {
  question: string;
  subjectId: number;
}

export interface Question {
  id: number;
  question: string;
  likesCount: number;
}

export interface QuestionList {
  questionList: Question[];
}

export interface SubjectList {
  subjectList: Subject[];
}

export interface QuestionId {
  id: number;
}

function createBaseLikes(): Likes {
  return { userId: "", questionId: 0 };
}

export const Likes = {
  encode(message: Likes, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.questionId !== 0) {
      writer.uint32(16).int64(message.questionId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Likes {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLikes();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.questionId = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Likes {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      questionId: isSet(object.questionId) ? Number(object.questionId) : 0,
    };
  },

  toJSON(message: Likes): unknown {
    const obj: any = {};
    if (message.userId !== "") {
      obj.userId = message.userId;
    }
    if (message.questionId !== 0) {
      obj.questionId = Math.round(message.questionId);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Likes>, I>>(base?: I): Likes {
    return Likes.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Likes>, I>>(object: I): Likes {
    const message = createBaseLikes();
    message.userId = object.userId ?? "";
    message.questionId = object.questionId ?? 0;
    return message;
  },
};

function createBaseNewSubject(): NewSubject {
  return { title: "" };
}

export const NewSubject = {
  encode(message: NewSubject, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NewSubject {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNewSubject();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.title = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): NewSubject {
    return { title: isSet(object.title) ? String(object.title) : "" };
  },

  toJSON(message: NewSubject): unknown {
    const obj: any = {};
    if (message.title !== "") {
      obj.title = message.title;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<NewSubject>, I>>(base?: I): NewSubject {
    return NewSubject.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<NewSubject>, I>>(object: I): NewSubject {
    const message = createBaseNewSubject();
    message.title = object.title ?? "";
    return message;
  },
};

function createBaseSubject(): Subject {
  return { id: 0, title: "" };
}

export const Subject = {
  encode(message: Subject, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).int64(message.id);
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Subject {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubject();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = longToNumber(reader.int64() as Long);
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.title = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Subject {
    return { id: isSet(object.id) ? Number(object.id) : 0, title: isSet(object.title) ? String(object.title) : "" };
  },

  toJSON(message: Subject): unknown {
    const obj: any = {};
    if (message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.title !== "") {
      obj.title = message.title;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Subject>, I>>(base?: I): Subject {
    return Subject.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Subject>, I>>(object: I): Subject {
    const message = createBaseSubject();
    message.id = object.id ?? 0;
    message.title = object.title ?? "";
    return message;
  },
};

function createBaseSubjectId(): SubjectId {
  return { id: 0 };
}

export const SubjectId = {
  encode(message: SubjectId, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).int64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubjectId {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubjectId();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SubjectId {
    return { id: isSet(object.id) ? Number(object.id) : 0 };
  },

  toJSON(message: SubjectId): unknown {
    const obj: any = {};
    if (message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SubjectId>, I>>(base?: I): SubjectId {
    return SubjectId.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<SubjectId>, I>>(object: I): SubjectId {
    const message = createBaseSubjectId();
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseNewQuestion(): NewQuestion {
  return { question: "", subjectId: 0 };
}

export const NewQuestion = {
  encode(message: NewQuestion, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.question !== "") {
      writer.uint32(10).string(message.question);
    }
    if (message.subjectId !== 0) {
      writer.uint32(16).int64(message.subjectId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NewQuestion {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNewQuestion();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.question = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.subjectId = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): NewQuestion {
    return {
      question: isSet(object.question) ? String(object.question) : "",
      subjectId: isSet(object.subjectId) ? Number(object.subjectId) : 0,
    };
  },

  toJSON(message: NewQuestion): unknown {
    const obj: any = {};
    if (message.question !== "") {
      obj.question = message.question;
    }
    if (message.subjectId !== 0) {
      obj.subjectId = Math.round(message.subjectId);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<NewQuestion>, I>>(base?: I): NewQuestion {
    return NewQuestion.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<NewQuestion>, I>>(object: I): NewQuestion {
    const message = createBaseNewQuestion();
    message.question = object.question ?? "";
    message.subjectId = object.subjectId ?? 0;
    return message;
  },
};

function createBaseQuestion(): Question {
  return { id: 0, question: "", likesCount: 0 };
}

export const Question = {
  encode(message: Question, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).int64(message.id);
    }
    if (message.question !== "") {
      writer.uint32(18).string(message.question);
    }
    if (message.likesCount !== 0) {
      writer.uint32(24).int64(message.likesCount);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Question {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuestion();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = longToNumber(reader.int64() as Long);
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.question = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.likesCount = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Question {
    return {
      id: isSet(object.id) ? Number(object.id) : 0,
      question: isSet(object.question) ? String(object.question) : "",
      likesCount: isSet(object.likesCount) ? Number(object.likesCount) : 0,
    };
  },

  toJSON(message: Question): unknown {
    const obj: any = {};
    if (message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.question !== "") {
      obj.question = message.question;
    }
    if (message.likesCount !== 0) {
      obj.likesCount = Math.round(message.likesCount);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Question>, I>>(base?: I): Question {
    return Question.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Question>, I>>(object: I): Question {
    const message = createBaseQuestion();
    message.id = object.id ?? 0;
    message.question = object.question ?? "";
    message.likesCount = object.likesCount ?? 0;
    return message;
  },
};

function createBaseQuestionList(): QuestionList {
  return { questionList: [] };
}

export const QuestionList = {
  encode(message: QuestionList, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.questionList) {
      Question.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuestionList {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuestionList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.questionList.push(Question.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): QuestionList {
    return {
      questionList: Array.isArray(object?.questionList)
        ? object.questionList.map((e: any) => Question.fromJSON(e))
        : [],
    };
  },

  toJSON(message: QuestionList): unknown {
    const obj: any = {};
    if (message.questionList?.length) {
      obj.questionList = message.questionList.map((e) => Question.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<QuestionList>, I>>(base?: I): QuestionList {
    return QuestionList.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<QuestionList>, I>>(object: I): QuestionList {
    const message = createBaseQuestionList();
    message.questionList = object.questionList?.map((e) => Question.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSubjectList(): SubjectList {
  return { subjectList: [] };
}

export const SubjectList = {
  encode(message: SubjectList, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.subjectList) {
      Subject.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubjectList {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubjectList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.subjectList.push(Subject.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SubjectList {
    return {
      subjectList: Array.isArray(object?.subjectList) ? object.subjectList.map((e: any) => Subject.fromJSON(e)) : [],
    };
  },

  toJSON(message: SubjectList): unknown {
    const obj: any = {};
    if (message.subjectList?.length) {
      obj.subjectList = message.subjectList.map((e) => Subject.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SubjectList>, I>>(base?: I): SubjectList {
    return SubjectList.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<SubjectList>, I>>(object: I): SubjectList {
    const message = createBaseSubjectList();
    message.subjectList = object.subjectList?.map((e) => Subject.fromPartial(e)) || [];
    return message;
  },
};

function createBaseQuestionId(): QuestionId {
  return { id: 0 };
}

export const QuestionId = {
  encode(message: QuestionId, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).int64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuestionId {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuestionId();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): QuestionId {
    return { id: isSet(object.id) ? Number(object.id) : 0 };
  },

  toJSON(message: QuestionId): unknown {
    const obj: any = {};
    if (message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<QuestionId>, I>>(base?: I): QuestionId {
    return QuestionId.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<QuestionId>, I>>(object: I): QuestionId {
    const message = createBaseQuestionId();
    message.id = object.id ?? 0;
    return message;
  },
};

export type BoardService = typeof BoardService;
export const BoardService = {
  postSubject: {
    path: "/board.Board/PostSubject",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: NewSubject) => Buffer.from(NewSubject.encode(value).finish()),
    requestDeserialize: (value: Buffer) => NewSubject.decode(value),
    responseSerialize: (value: Subject) => Buffer.from(Subject.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Subject.decode(value),
  },
  deleteSubject: {
    path: "/board.Board/DeleteSubject",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: SubjectId) => Buffer.from(SubjectId.encode(value).finish()),
    requestDeserialize: (value: Buffer) => SubjectId.decode(value),
    responseSerialize: (value: Empty) => Buffer.from(Empty.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Empty.decode(value),
  },
  listSubject: {
    path: "/board.Board/ListSubject",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: Empty) => Buffer.from(Empty.encode(value).finish()),
    requestDeserialize: (value: Buffer) => Empty.decode(value),
    responseSerialize: (value: SubjectList) => Buffer.from(SubjectList.encode(value).finish()),
    responseDeserialize: (value: Buffer) => SubjectList.decode(value),
  },
  postQuestion: {
    path: "/board.Board/PostQuestion",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: NewQuestion) => Buffer.from(NewQuestion.encode(value).finish()),
    requestDeserialize: (value: Buffer) => NewQuestion.decode(value),
    responseSerialize: (value: Question) => Buffer.from(Question.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Question.decode(value),
  },
  deleteQuestion: {
    path: "/board.Board/DeleteQuestion",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: QuestionId) => Buffer.from(QuestionId.encode(value).finish()),
    requestDeserialize: (value: Buffer) => QuestionId.decode(value),
    responseSerialize: (value: Empty) => Buffer.from(Empty.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Empty.decode(value),
  },
  listQuestion: {
    path: "/board.Board/ListQuestion",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: SubjectId) => Buffer.from(SubjectId.encode(value).finish()),
    requestDeserialize: (value: Buffer) => SubjectId.decode(value),
    responseSerialize: (value: QuestionList) => Buffer.from(QuestionList.encode(value).finish()),
    responseDeserialize: (value: Buffer) => QuestionList.decode(value),
  },
  getQuestion: {
    path: "/board.Board/GetQuestion",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: QuestionId) => Buffer.from(QuestionId.encode(value).finish()),
    requestDeserialize: (value: Buffer) => QuestionId.decode(value),
    responseSerialize: (value: Question) => Buffer.from(Question.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Question.decode(value),
  },
  postLikes: {
    path: "/board.Board/PostLikes",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: Likes) => Buffer.from(Likes.encode(value).finish()),
    requestDeserialize: (value: Buffer) => Likes.decode(value),
    responseSerialize: (value: Empty) => Buffer.from(Empty.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Empty.decode(value),
  },
  deleteLikes: {
    path: "/board.Board/DeleteLikes",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: Likes) => Buffer.from(Likes.encode(value).finish()),
    requestDeserialize: (value: Buffer) => Likes.decode(value),
    responseSerialize: (value: Empty) => Buffer.from(Empty.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Empty.decode(value),
  },
} as const;

export interface BoardServer extends UntypedServiceImplementation {
  postSubject: handleUnaryCall<NewSubject, Subject>;
  deleteSubject: handleUnaryCall<SubjectId, Empty>;
  listSubject: handleUnaryCall<Empty, SubjectList>;
  postQuestion: handleUnaryCall<NewQuestion, Question>;
  deleteQuestion: handleUnaryCall<QuestionId, Empty>;
  listQuestion: handleUnaryCall<SubjectId, QuestionList>;
  getQuestion: handleUnaryCall<QuestionId, Question>;
  postLikes: handleUnaryCall<Likes, Empty>;
  deleteLikes: handleUnaryCall<Likes, Empty>;
}

export interface BoardClient extends Client {
  postSubject(request: NewSubject, callback: (error: ServiceError | null, response: Subject) => void): ClientUnaryCall;
  postSubject(
    request: NewSubject,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: Subject) => void,
  ): ClientUnaryCall;
  postSubject(
    request: NewSubject,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: Subject) => void,
  ): ClientUnaryCall;
  deleteSubject(request: SubjectId, callback: (error: ServiceError | null, response: Empty) => void): ClientUnaryCall;
  deleteSubject(
    request: SubjectId,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: Empty) => void,
  ): ClientUnaryCall;
  deleteSubject(
    request: SubjectId,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: Empty) => void,
  ): ClientUnaryCall;
  listSubject(request: Empty, callback: (error: ServiceError | null, response: SubjectList) => void): ClientUnaryCall;
  listSubject(
    request: Empty,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: SubjectList) => void,
  ): ClientUnaryCall;
  listSubject(
    request: Empty,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: SubjectList) => void,
  ): ClientUnaryCall;
  postQuestion(
    request: NewQuestion,
    callback: (error: ServiceError | null, response: Question) => void,
  ): ClientUnaryCall;
  postQuestion(
    request: NewQuestion,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: Question) => void,
  ): ClientUnaryCall;
  postQuestion(
    request: NewQuestion,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: Question) => void,
  ): ClientUnaryCall;
  deleteQuestion(request: QuestionId, callback: (error: ServiceError | null, response: Empty) => void): ClientUnaryCall;
  deleteQuestion(
    request: QuestionId,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: Empty) => void,
  ): ClientUnaryCall;
  deleteQuestion(
    request: QuestionId,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: Empty) => void,
  ): ClientUnaryCall;
  listQuestion(
    request: SubjectId,
    callback: (error: ServiceError | null, response: QuestionList) => void,
  ): ClientUnaryCall;
  listQuestion(
    request: SubjectId,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: QuestionList) => void,
  ): ClientUnaryCall;
  listQuestion(
    request: SubjectId,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: QuestionList) => void,
  ): ClientUnaryCall;
  getQuestion(request: QuestionId, callback: (error: ServiceError | null, response: Question) => void): ClientUnaryCall;
  getQuestion(
    request: QuestionId,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: Question) => void,
  ): ClientUnaryCall;
  getQuestion(
    request: QuestionId,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: Question) => void,
  ): ClientUnaryCall;
  postLikes(request: Likes, callback: (error: ServiceError | null, response: Empty) => void): ClientUnaryCall;
  postLikes(
    request: Likes,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: Empty) => void,
  ): ClientUnaryCall;
  postLikes(
    request: Likes,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: Empty) => void,
  ): ClientUnaryCall;
  deleteLikes(request: Likes, callback: (error: ServiceError | null, response: Empty) => void): ClientUnaryCall;
  deleteLikes(
    request: Likes,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: Empty) => void,
  ): ClientUnaryCall;
  deleteLikes(
    request: Likes,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: Empty) => void,
  ): ClientUnaryCall;
}

export const BoardClient = makeGenericClientConstructor(BoardService, "board.Board") as unknown as {
  new (address: string, credentials: ChannelCredentials, options?: Partial<ClientOptions>): BoardClient;
  service: typeof BoardService;
};

declare const self: any | undefined;
declare const window: any | undefined;
declare const global: any | undefined;
const tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
